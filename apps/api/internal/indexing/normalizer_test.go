package indexing

import (
	"testing"
)

func TestNormalizeText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "crlf to lf",
			input:    "line1\r\nline2\r\n",
			expected: "line1\nline2",
		},
		{
			name:     "trailing whitespace trimmed outside code blocks",
			input:    "hello world   \nfoo bar\t\t\n",
			expected: "hello world\nfoo bar",
		},
		{
			name:     "excessive blank lines collapsed outside code blocks",
			input:    "para1\n\n\n\npara2\n\n\npara3",
			expected: "para1\n\npara2\n\npara3",
		},
		{
			name: "verbatim fenced code block preserved exactly",
			input: "Intro text    \n\n```go\nfunc main() {  \n    \n    println(\"hello\")   \n}\n```\n\nOutro text",
			expected: "Intro text\n\n```go\nfunc main() {  \n    \n    println(\"hello\")   \n}\n```\n\nOutro text",
		},
		{
			name:     "empty input",
			input:    "   \n\r\n\t  ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NormalizeText(tt.input)
			if actual != tt.expected {
				t.Errorf("NormalizeText() mismatch\nExpected:\n%q\nActual:\n%q", tt.expected, actual)
			}
		})
	}
}
