package parser

import (
	"io"
	"strings"
)

// DOCXParser handles DOCX file parsing
type DOCXParser struct{}

// NewDOCXParser creates a new DOCX parser
func NewDOCXParser() *DOCXParser {
	return &DOCXParser{}
}

// Parse extracts text from DOCX file
// TODO: Implement actual DOCX parsing using a library like unidoc/unioffice
func (p *DOCXParser) Parse(reader io.Reader) (string, error) {
	// Placeholder implementation
	// In production, use: github.com/unidoc/unioffice

	// For MVP, return a simple message
	var builder strings.Builder
	builder.WriteString("[DOCX Content]\n")
	builder.WriteString("TODO: Implement DOCX text extraction\n")
	builder.WriteString("Use library like unidoc/unioffice\n")

	return builder.String(), nil
}

// SupportedType returns the file type this parser supports
func (p *DOCXParser) SupportedType() string {
	return "docx"
}
