package indexing

import (
	"fmt"
	"strings"
	"time"

	"github.com/nachsyas/arham-porto/apps/api/internal/domain"
)

const (
	MaxChunkCharSize = 1500
	ChunkOverlapSize = 200
)

// ChunkDocument segments a normalized source document into deterministic KnowledgeChunks.
// Prepends deterministic repository, path, and section context per Gate #32.
func ChunkDocument(
	source domain.KnowledgeSource,
	normalizedContent string,
	projectID *string,
	skillIDs []string,
	evidenceID *string,
) []domain.KnowledgeChunk {
	if strings.TrimSpace(normalizedContent) == "" {
		return nil
	}

	isMarkdown := strings.HasSuffix(strings.ToLower(source.Path), ".md") || source.SourceType == "doc"

	var rawChunks []rawChunkData
	if isMarkdown {
		rawChunks = chunkMarkdown(source, normalizedContent)
	} else {
		rawChunks = chunkCodeOrText(source, normalizedContent)
	}

	var chunks []domain.KnowledgeChunk
	for idx, rc := range rawChunks {
		// Context prepending per Gate #32
		contextHeader := fmt.Sprintf("Repository: %s\nPath: %s", source.Repository, source.Path)
		if rc.section != "" {
			contextHeader += fmt.Sprintf("\nSection: %s", rc.section)
		}
		fullContent := fmt.Sprintf("%s\n\n%s", contextHeader, strings.TrimSpace(rc.content))

		checksum := CalculateStringSHA256(fullContent)
		chunkID := GenerateChunkID(source.ID, idx, checksum)

		chunks = append(chunks, domain.KnowledgeChunk{
			ID:                  chunkID,
			SourceID:            source.ID,
			ChunkIndex:          idx,
			Content:             fullContent,
			Checksum:            checksum,
			ProjectID:           projectID,
			SkillIDs:            skillIDs,
			EvidenceID:          evidenceID,
			CreatedAt:           time.Now().UTC(),
		})
	}

	return chunks
}

type rawChunkData struct {
	section string
	content string
}

func chunkMarkdown(source domain.KnowledgeSource, content string) []rawChunkData {
	lines := strings.Split(content, "\n")
	var chunks []rawChunkData

	currentSection := "Overview"
	var currentLines []string
	currentLen := 0

	flush := func() {
		if len(currentLines) == 0 {
			return
		}
		text := strings.Join(currentLines, "\n")
		trimmed := strings.TrimSpace(text)
		if trimmed != "" {
			if len(trimmed) > MaxChunkCharSize {
				subChunks := splitLongText(trimmed, MaxChunkCharSize, ChunkOverlapSize)
				for subIdx, sc := range subChunks {
					sec := currentSection
					if len(subChunks) > 1 {
						sec = fmt.Sprintf("%s (Part %d)", currentSection, subIdx+1)
					}
					chunks = append(chunks, rawChunkData{section: sec, content: sc})
				}
			} else {
				chunks = append(chunks, rawChunkData{section: currentSection, content: trimmed})
			}
		}
		currentLines = nil
		currentLen = 0
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			// Check if it's a markdown header
			parts := strings.SplitN(trimmed, " ", 2)
			if len(parts) == 2 && strings.Count(parts[0], "#") == len(parts[0]) {
				flush()
				currentSection = strings.TrimSpace(parts[1])
				currentLines = append(currentLines, line)
				currentLen += len(line) + 1
				continue
			}
		}

		currentLines = append(currentLines, line)
		currentLen += len(line) + 1

		if currentLen >= MaxChunkCharSize {
			flush()
		}
	}

	flush()
	return chunks
}

func chunkCodeOrText(source domain.KnowledgeSource, content string) []rawChunkData {
	lines := strings.Split(content, "\n")
	var chunks []rawChunkData

	if len(lines) == 0 {
		return nil
	}

	const linesPerChunk = 40
	const overlapLines = 5

	startLine := 0
	for startLine < len(lines) {
		endLine := startLine + linesPerChunk
		if endLine > len(lines) {
			endLine = len(lines)
		}

		chunkLines := lines[startLine:endLine]
		text := strings.TrimSpace(strings.Join(chunkLines, "\n"))
		if text != "" {
			section := fmt.Sprintf("Lines %d-%d", startLine+1, endLine)
			chunks = append(chunks, rawChunkData{section: section, content: text})
		}

		if endLine == len(lines) {
			break
		}
		startLine += (linesPerChunk - overlapLines)
	}

	return chunks
}

func splitLongText(text string, maxLen, overlap int) []string {
	var result []string
	paras := strings.Split(text, "\n\n")

	var currentParas []string
	currentLen := 0

	for _, p := range paras {
		pTrim := strings.TrimSpace(p)
		if pTrim == "" {
			continue
		}

		if currentLen+len(pTrim) > maxLen && len(currentParas) > 0 {
			result = append(result, strings.Join(currentParas, "\n\n"))
			currentParas = nil
			currentLen = 0
		}

		currentParas = append(currentParas, pTrim)
		currentLen += len(pTrim) + 2
	}

	if len(currentParas) > 0 {
		result = append(result, strings.Join(currentParas, "\n\n"))
	}

	return result
}
