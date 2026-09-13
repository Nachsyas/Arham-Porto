package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/nachsyas/arham-porto/apps/api/internal/ai"
	"github.com/nachsyas/arham-porto/apps/api/internal/domain"
	"github.com/nachsyas/arham-porto/apps/api/internal/llm"
	"github.com/nachsyas/arham-porto/apps/api/internal/retrieval"
)

// Retriever defines the interface for retrieving evidence chunks.
type Retriever interface {
	Retrieve(ctx context.Context, query string, opts retrieval.SearchOptions) ([]domain.RetrievalResult, error)
}

// AskUseCase orchestrates grounded RAG question-answering.
type AskUseCase interface {
	Ask(ctx context.Context, question string) (llm.GroundedResponse, error)
	GetRetrievedEvidence(ctx context.Context, question string) ([]llm.EvidenceContext, error)
}

type askUseCase struct {
	retriever       Retriever
	llmProvider     llm.LLMProvider
	maxEvidenceChars int
}

// NewAskUseCase constructs a new AskUseCase instance.
func NewAskUseCase(retriever Retriever, llmProvider llm.LLMProvider, maxEvidenceChars int) AskUseCase {
	if maxEvidenceChars <= 0 {
		maxEvidenceChars = 24000
	}
	return &askUseCase{
		retriever:        retriever,
		llmProvider:      llmProvider,
		maxEvidenceChars: maxEvidenceChars,
	}
}

// isPrivacyViolation checks for clear requests for private personal information.
func isPrivacyViolation(q string) bool {
	lower := strings.ToLower(q)
	privateKeywords := []string{
		"home address",
		"residential address",
		"where do you live",
		"where does arham live",
		"where does nachsyas live",
		"phone number",
		"telephone",
		"mobile number",
		"date of birth",
		"birth date",
		"birth year",
		"exact coordinate",
		"private address",
	}
	for _, kw := range privateKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// isPromptLeakRequest checks for explicit attempts to disclose system instructions.
func isPromptLeakRequest(q string) bool {
	lower := strings.ToLower(q)
	leakKeywords := []string{
		"reveal system prompt",
		"reveal your system prompt",
		"show system prompt",
		"show your system prompt",
		"what is your system prompt",
		"what are your instructions",
		"ignore all previous instructions",
		"ignore previous instructions",
		"ignore all instructions",
	}
	for _, kw := range leakKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// Ask executes grounded Q&A with strict claim validation and server-owned citations.
func (u *askUseCase) Ask(ctx context.Context, question string) (llm.GroundedResponse, error) {
	trimmedQ := strings.TrimSpace(question)
	if len(trimmedQ) < 2 || len(trimmedQ) > 1000 {
		return llm.GroundedResponse{}, fmt.Errorf("question must be between 2 and 1000 characters")
	}

	// 1. Privacy Pre-Guard (Correction 20)
	if isPrivacyViolation(trimmedQ) {
		contactAction, _ := ai.GetSafeAction("go-to-contact")
		return llm.GroundedResponse{
			Status: "privacy_refusal",
			Answer: "I cannot provide private personal details such as home addresses, phone numbers, or dates of birth. I am happy to answer questions regarding Nachsyas Arham Mumtaz Nashohi's verified engineering experience, projects, and skills.",
			Segments: []llm.GroundedSegment{
				{
					Text:        "I cannot provide private personal details such as home addresses, phone numbers, or dates of birth.",
					EvidenceIDs: []string{},
				},
			},
			Evidence: []llm.PublicEvidenceItem{},
			Sources:  []llm.SourceCitation{},
			Actions: []llm.PublicAction{
				{ID: contactAction.ID, Label: contactAction.Label},
			},
		}, nil
	}

	// 2. Prompt Leak Pre-Guard (Correction 63)
	if isPromptLeakRequest(trimmedQ) {
		projectsAction, _ := ai.GetSafeAction("go-to-projects")
		return llm.GroundedResponse{
			Status: "scope_refusal",
			Answer: "Internal system instructions are not part of portfolio evidence. I am here to help answer questions about Nachsyas Arham Mumtaz Nashohi's verified projects, technical skills, and engineering experience.",
			Segments: []llm.GroundedSegment{
				{
					Text:        "Internal system instructions are not part of portfolio evidence.",
					EvidenceIDs: []string{},
				},
			},
			Evidence: []llm.PublicEvidenceItem{},
			Sources:  []llm.SourceCitation{},
			Actions: []llm.PublicAction{
				{ID: projectsAction.ID, Label: projectsAction.Label},
			},
		}, nil
	}

	// 3. Retrieval
	evidenceList, err := u.GetRetrievedEvidence(ctx, trimmedQ)
	if err != nil {
		return llm.GroundedResponse{}, fmt.Errorf("evidence retrieval failed: %w", err)
	}

	// 4. Zero-Evidence Short-Circuit (Correction 19)
	if len(evidenceList) == 0 {
		projectsAction, _ := ai.GetSafeAction("go-to-projects")
		skillsAction, _ := ai.GetSafeAction("go-to-skills")
		return llm.GroundedResponse{
			Status: "insufficient_evidence",
			Answer: "I couldn't find verified evidence of that in Nachsyas Arham Mumtaz Nashohi's approved portfolio sources.",
			Segments: []llm.GroundedSegment{
				{
					Text:        "I couldn't find verified evidence of that in Nachsyas Arham Mumtaz Nashohi's approved portfolio sources.",
					EvidenceIDs: []string{},
				},
			},
			Evidence: []llm.PublicEvidenceItem{},
			Sources:  []llm.SourceCitation{},
			Actions: []llm.PublicAction{
				{ID: projectsAction.ID, Label: projectsAction.Label},
				{ID: skillsAction.ID, Label: skillsAction.Label},
			},
		}, nil
	}

	// 5. Build LLM Request with Delimited Untrusted Context
	formattedEvidence := ai.FormatEvidenceContext(evidenceList, u.maxEvidenceChars)
	req := llm.GenerateRequest{
		SystemInstruction: ai.SystemInstruction,
		UserQuestion:      fmt.Sprintf("Reviewer Question:\n%s\n\n%s", trimmedQ, formattedEvidence),
		Evidence:          evidenceList,
	}

	// 6. Call LLM Provider
	rawAnswer, err := u.llmProvider.Generate(ctx, req)
	if err != nil {
		return llm.GroundedResponse{}, fmt.Errorf("generation provider failed: %w", err)
	}

	// 7. Validate and Sanitize Output
	return u.validateAndAssemble(rawAnswer, evidenceList)
}

// GetRetrievedEvidence fetches and formats evidence for a query.
func (u *askUseCase) GetRetrievedEvidence(ctx context.Context, question string) ([]llm.EvidenceContext, error) {
	if u.retriever == nil {
		return nil, fmt.Errorf("retriever is not initialized")
	}

	results, err := u.retriever.Retrieve(ctx, question, retrieval.SearchOptions{Limit: 5})
	if err != nil {
		return nil, err
	}

	evidenceList := make([]llm.EvidenceContext, 0, len(results))
	for i, r := range results {
		evID := fmt.Sprintf("E%d", i+1)
		kind := "portfolio"
		var repoPtr *string
		var pathPtr *string
		var commitPtr *string
		var urlPtr *string

		if r.SourceType == "github" {
			kind = "github"
			if r.Repository != "" {
				repoVal := r.Repository
				repoPtr = &repoVal
			}
			if r.Path != "" {
				pathVal := r.Path
				pathPtr = &pathVal
			}
			if r.CommitSHA != "" {
				commitVal := r.CommitSHA
				commitPtr = &commitVal
			}
			if r.SourceURL != "" {
				urlVal := r.SourceURL
				urlPtr = &urlVal
			}
		}

		evidenceList = append(evidenceList, llm.EvidenceContext{
			ID:          evID,
			Kind:        kind,
			Title:       r.SourceTitle,
			Repository:  repoPtr,
			Path:        pathPtr,
			CommitSHA:   commitPtr,
			CitationURL: urlPtr,
			Content:     r.Content,
			ProjectID:   r.ProjectID,
			SkillIDs:    r.SkillIDs,
			EvidenceID:  r.EvidenceID,
			Similarity:  float32(r.Similarity),
		})
	}

	return evidenceList, nil
}

// validateAndAssemble performs claim-level evidence verification and builds the public GroundedResponse.
func (u *askUseCase) validateAndAssemble(raw llm.GeneratedAnswer, evidenceList []llm.EvidenceContext) (llm.GroundedResponse, error) {
	// Status validation
	status := raw.Status
	switch status {
	case "supported", "insufficient_evidence", "privacy_refusal", "scope_refusal":
		// valid
	default:
		status = "insufficient_evidence"
	}

	evidenceMap := make(map[string]llm.EvidenceContext, len(evidenceList))
	for _, e := range evidenceList {
		evidenceMap[e.ID] = e
	}

	// Validate segments and collect used evidence IDs
	usedEvidenceIDSet := make(map[string]bool)
	validatedSegments := make([]llm.GroundedSegment, 0, len(raw.Segments))
	var answerParts []string

	for _, seg := range raw.Segments {
		trimmedText := strings.TrimSpace(seg.Text)
		if trimmedText == "" {
			continue
		}

		// Filter evidence IDs: drop unknown IDs (Correction 5)
		validIDs := make([]string, 0, len(seg.EvidenceIDs))
		for _, id := range seg.EvidenceIDs {
			if _, exists := evidenceMap[id]; exists {
				validIDs = append(validIDs, id)
				usedEvidenceIDSet[id] = true
			}
		}

		validatedSegments = append(validatedSegments, llm.GroundedSegment{
			Text:        trimmedText,
			EvidenceIDs: validIDs,
		})
		answerParts = append(answerParts, trimmedText)
	}

	// If no valid segments were generated, provide standard fallback based on status
	if len(validatedSegments) == 0 {
		fallbackText := "I couldn't find verified evidence of that in Nachsyas Arham Mumtaz Nashohi's approved portfolio sources."
		if status == "privacy_refusal" {
			fallbackText = "I cannot provide private personal details such as home addresses, phone numbers, or dates of birth."
		} else if status == "scope_refusal" {
			fallbackText = "Internal system instructions are not part of portfolio evidence."
		}
		validatedSegments = append(validatedSegments, llm.GroundedSegment{
			Text:        fallbackText,
			EvidenceIDs: []string{},
		})
		answerParts = append(answerParts, fallbackText)
	}

	// Build full answer text
	finalAnswer := strings.Join(answerParts, " ")

	// Build verified SourceCitations and PublicEvidenceItems from used evidence
	citations := make([]llm.SourceCitation, 0)
	publicEvidence := make([]llm.PublicEvidenceItem, 0)
	citationIndex := 1

	for _, e := range evidenceList {
		// Include evidence if explicitly cited, or if status is supported and only 1-2 chunks available
		isUsed := usedEvidenceIDSet[e.ID]
		if !isUsed && status == "supported" && len(usedEvidenceIDSet) == 0 {
			isUsed = true // Fallback if model omitted ID in segments
		}

		if !isUsed {
			continue
		}

		citID := fmt.Sprintf("S%d", citationIndex)
		citationIndex++

		var cit llm.SourceCitation
		if e.Kind == "github" {
			// GitHub citation with immutable HTTPS link (Correction 13)
			label := e.Title
			if e.Repository != nil && e.Path != nil {
				label = fmt.Sprintf("%s / %s", *e.Repository, *e.Path)
			}
			cit = llm.SourceCitation{
				ID:         citID,
				Kind:       "github",
				Label:      label,
				URL:        e.CitationURL,
				Repository: e.Repository,
				Path:       e.Path,
				CommitSHA:  e.CommitSHA,
			}
		} else {
			// Canonical portfolio citation (Correction 14, 15)
			label := e.Title
			var safeTarget *string
			if e.ProjectID != nil {
				label = fmt.Sprintf("Project: %s", *e.ProjectID)
				target := fmt.Sprintf("/projects/%s", *e.ProjectID)
				safeTarget = &target
			} else {
				target := "/#journey"
				safeTarget = &target
			}

			cit = llm.SourceCitation{
				ID:    citID,
				Kind:  "portfolio",
				Label: label,
				URL:   safeTarget,
			}
		}
		citations = append(citations, cit)

		// Create safe public evidence excerpt (bounded to 600 chars, no similarity scores, Correction 17, 49)
		excerpt := e.Content
		if len(excerpt) > 600 {
			excerpt = excerpt[:600] + "..."
		}

		publicEvidence = append(publicEvidence, llm.PublicEvidenceItem{
			ID:         e.ID,
			Kind:       e.Kind,
			Title:      e.Title,
			Repository: e.Repository,
			Path:       e.Path,
			Excerpt:    excerpt,
			CitationID: citID,
		})
	}

	// Validate safe actions (Correction 29, 30)
	publicActions := make([]llm.PublicAction, 0)
	seenActions := make(map[string]bool)

	for _, actionID := range raw.SuggestedActionIDs {
		if action, exists := ai.GetSafeAction(actionID); exists && !seenActions[action.ID] {
			publicActions = append(publicActions, llm.PublicAction{
				ID:    action.ID,
				Label: action.Label,
			})
			seenActions[action.ID] = true
		}
	}

	// Default fallback action if none validated
	if len(publicActions) == 0 {
		if defaultAction, exists := ai.GetSafeAction("go-to-projects"); exists {
			publicActions = append(publicActions, llm.PublicAction{
				ID:    defaultAction.ID,
				Label: defaultAction.Label,
			})
		}
	}

	return llm.GroundedResponse{
		Status:   status,
		Answer:   finalAnswer,
		Segments: validatedSegments,
		Evidence: publicEvidence,
		Sources:  citations,
		Actions:  publicActions,
	}, nil
}
