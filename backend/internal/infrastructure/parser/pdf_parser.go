package parser

import (
	"io"
	"strings"
)

// PDFParser handles PDF file parsing
type PDFParser struct{}

// NewPDFParser creates a new PDF parser
func NewPDFParser() *PDFParser {
	return &PDFParser{}
}

// Parse extracts text from PDF file
// TODO: Implement actual PDF parsing using a library like unidoc or pdfcpu
func (p *PDFParser) Parse(reader io.Reader) (string, error) {
	// Placeholder implementation
	// In production, use: github.com/ledongthuc/pdf or github.com/unidoc/unipdf

	// For MVP, return a simple message
	// Real implementation would extract text from PDF
	var builder strings.Builder
	builder.WriteString("[PDF Content]\n")
	builder.WriteString("TODO: Implement PDF text extraction\n")
	builder.WriteString("Use library like unidoc/unipdf or ledongthuc/pdf\n")

	return builder.String(), nil
}

// SupportedType returns the file type this parser supports
func (p *PDFParser) SupportedType() string {
	return "pdf"
}
