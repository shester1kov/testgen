package parser

import (
	"io"
)

// MDParser implements DocumentParser for Markdown files
type MDParser struct{}

// NewMDParser creates a new Markdown parser
func NewMDParser() *MDParser {
	return &MDParser{}
}

// Parse extracts text from a Markdown file
func (p *MDParser) Parse(reader io.Reader) (string, error) {
	// Markdown is plain text, so we can read it directly
	content, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// SupportedType returns the file type this parser supports
func (p *MDParser) SupportedType() string {
	return "md"
}
