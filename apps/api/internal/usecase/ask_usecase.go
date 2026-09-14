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
	retriever        Retriever
	llmProvider      llm.LLMProvider
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

// GetRetrievedEvidence fetches and formats evidence for a query with correct source family classification.
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

		// Source family classification per Correction 6 & 7:
		// repository == "canonical" -> portfolio
		// repository != "canonical" -> github (source types: doc, manifest, architecture, entrypoint)
		isGitHub := r.Repository != "canonical" && r.Repository != ""
		kind := "portfolio"
		var repoPtr *string
		var pathPtr *string
		var commitPtr *string
		var urlPtr *string

		if isGitHub {
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
			if r.SourceURL != "" && strings.HasPrefix(r.SourceURL, "https://") {
				urlVal := r.SourceURL
				urlPtr = &urlVal
			} else if repoPtr != nil && pathPtr != nil && commitPtr != nil && len(*commitPtr) == 40 {
				urlVal := fmt.Sprintf("https://github.com/%s/blob/%s/%s", *repoPtr, *commitPtr, *pathPtr)
				urlPtr = &urlVal
			}
		}

		evidenceList = append(evidenceList, llm.EvidenceContext{
			ID:          evID,
			Kind:        kind,
			SourceType:  r.SourceType,
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

	// Validate segments and collect used evidence IDs (Corrections 9, 10, 11, 12)
	usedEvidenceIDSet := make(map[string]bool)
	validatedSegments := make([]llm.GroundedSegment, 0, len(raw.Segments))
	var answerParts []string

	for _, seg := range raw.Segments {
		trimmedText := strings.TrimSpace(seg.Text)
		if trimmedText == "" {
			continue
		}

		// Filter evidence IDs: drop unknown IDs (Correction 5, 11, 12)
		validIDs := make([]string, 0, len(seg.EvidenceIDs))
		for _, id := range seg.EvidenceIDs {
			if _, exists := evidenceMap[id]; exists {
				validIDs = append(validIDs, id)
			}
		}

		// Correction 9: For status == "supported", EVERY non-empty factual answer segment
		// must contain at least ONE valid retrieved evidence ID after server validation.
		// If zero valid evidence IDs remain, DO NOT retain that segment.
		if status == "supported" && len(validIDs) == 0 {
			continue
		}

		// Track used evidence IDs for retained segments
		for _, id := range validIDs {
			usedEvidenceIDSet[id] = true
		}

		validatedSegments = append(validatedSegments, llm.GroundedSegment{
			Text:        trimmedText,
			EvidenceIDs: validIDs,
		})
		answerParts = append(answerParts, trimmedText)
	}

	// Correction 9, 10, 11: If status is supported but no valid grounded segments remain,
	// downgrade entire result to insufficient_evidence with deterministic safe response.
	if status == "supported" && (len(validatedSegments) == 0 || len(usedEvidenceIDSet) == 0) {
		status = "insufficient_evidence"
		fallbackText := "I couldn't find verified evidence of that in Nachsyas Arham Mumtaz Nashohi's approved portfolio sources."
		validatedSegments = []llm.GroundedSegment{
			{
				Text:        fallbackText,
				EvidenceIDs: []string{},
			},
		}
		answerParts = []string{fallbackText}
		usedEvidenceIDSet = make(map[string]bool)
	}

	// For other statuses, if no segments were provided:
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

	// Build verified SourceCitations and PublicEvidenceItems from strictly cited evidence (Correction 10)
	citations := make([]llm.SourceCitation, 0)
	publicEvidence := make([]llm.PublicEvidenceItem, 0)
	citationIndex := 1

	for _, e := range evidenceList {
		// Include evidence ONLY if explicitly associated with a retained segment (Correction 10)
		if !usedEvidenceIDSet[e.ID] {
			continue
		}

		citID := fmt.Sprintf("S%d", citationIndex)
		citationIndex++

		var cit llm.SourceCitation
		if e.Kind == "github" {
			// GitHub citation with immutable HTTPS link (Correction 7, 13)
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
			// Canonical portfolio citation routing per Correction 8:
			// canonical_profile  -> label: "Portfolio Profile"   -> URL: nil / optional
			// canonical_project  -> label: "Project: <title>"    -> URL: /projects/<slug>
			// canonical_skill    -> label: "Skills & Evidence"   -> URL: /#skills
			// canonical_evidence -> label: "Skills & Evidence"   -> URL: /#skills (or /projects/<slug>)
			// canonical_journey  -> label: "Academic Journey"    -> URL: /#journey
			label := e.Title
			var safeTarget *string

			switch e.SourceType {
			case "canonical_profile":
				label = "Portfolio Profile"
				safeTarget = nil
			case "canonical_project":
				if e.ProjectID != nil && *e.ProjectID != "" {
					label = fmt.Sprintf("Project: %s", *e.ProjectID)
					target := fmt.Sprintf("/projects/%s", *e.ProjectID)
					safeTarget = &target
				} else {
					label = "Projects"
					target := "/#projects"
					safeTarget = &target
				}
			case "canonical_skill":
				label = "Skills & Evidence"
				target := "/#skills"
				safeTarget = &target
			case "canonical_evidence":
				label = "Skills & Evidence"
				if e.ProjectID != nil && *e.ProjectID != "" {
					target := fmt.Sprintf("/projects/%s", *e.ProjectID)
					safeTarget = &target
				} else {
					target := "/#skills"
					safeTarget = &target
				}
			case "canonical_journey":
				label = "Academic Journey"
				target := "/#journey"
				safeTarget = &target
			default:
				if e.ProjectID != nil && *e.ProjectID != "" {
					label = fmt.Sprintf("Project: %s", *e.ProjectID)
					target := fmt.Sprintf("/projects/%s", *e.ProjectID)
					safeTarget = &target
				} else {
					target := "/#projects"
					safeTarget = &target
				}
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

	// Status consistency check per Correction 13:
	// supported must have at least one grounded segment and valid citation
	if status == "supported" && (len(citations) == 0 || len(publicEvidence) == 0 || len(validatedSegments) == 0) {
		status = "insufficient_evidence"
		fallbackText := "I couldn't find verified evidence of that in Nachsyas Arham Mumtaz Nashohi's approved portfolio sources."
		validatedSegments = []llm.GroundedSegment{
			{
				Text:        fallbackText,
				EvidenceIDs: []string{},
			},
		}
		finalAnswer = fallbackText
		citations = []llm.SourceCitation{}
		publicEvidence = []llm.PublicEvidenceItem{}
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
