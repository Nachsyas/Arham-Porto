package indexing

import (
	"strings"
)

// NormalizeText cleans and standardizes text content for deterministic indexing.
// Rules per Gate #30:
// 1. Line endings are normalized globally (\r\n -> \n, \r -> \n).
// 2. Outside fenced code blocks:
//    - Trim trailing whitespace on each line
//    - Collapse excessive blank lines (maximum 1 blank line between text blocks)
// 3. Inside fenced code blocks:
//    - Preserve content verbatim (spaces, indentation, blank lines), only line-ending normalization applies.
func NormalizeText(content string) string {
	if len(content) == 0 {
		return ""
	}

	// 1. Globally normalize line endings
	normalizedNewlines := strings.ReplaceAll(content, "\r\n", "\n")
	normalizedNewlines = strings.ReplaceAll(normalizedNewlines, "\r", "\n")

	lines := strings.Split(normalizedNewlines, "\n")
	var result []string
	inCodeBlock := false
	fenceMarker := ""
	consecutiveBlankLines := 0

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Check for code fence start/end
		if !inCodeBlock {
			if strings.HasPrefix(trimmedLine, "```") || strings.HasPrefix(trimmedLine, "~~~") {
				inCodeBlock = true
				if strings.HasPrefix(trimmedLine, "```") {
					fenceMarker = "```"
				} else {
					fenceMarker = "~~~"
				}
				// Output fence line with trailing space trimmed
				result = append(result, strings.TrimRight(line, " \t"))
				consecutiveBlankLines = 0
				continue
			}
		} else {
			// Inside code block - check for closing fence
			if strings.HasPrefix(trimmedLine, fenceMarker) {
				inCodeBlock = false
				fenceMarker = ""
				result = append(result, strings.TrimRight(line, " \t"))
				consecutiveBlankLines = 0
				continue
			}
			// Inside code block - preserve line verbatim
			result = append(result, line)
			consecutiveBlankLines = 0
			continue
		}

		// Outside code block
		cleanedLine := strings.TrimRight(line, " \t")
		if len(cleanedLine) == 0 {
			consecutiveBlankLines++
			if consecutiveBlankLines <= 1 {
				result = append(result, "")
			}
			continue
		}

		consecutiveBlankLines = 0
		result = append(result, cleanedLine)
	}

	// Clean any leading or trailing blank lines
	start := 0
	for start < len(result) && result[start] == "" {
		start++
	}
	end := len(result)
	for end > start && result[end-1] == "" {
		end--
	}

	if start >= end {
		return ""
	}

	return strings.Join(result[start:end], "\n")
}
