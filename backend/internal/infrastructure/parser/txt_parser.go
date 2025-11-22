package parser

import (
	"io"
)

// TXTParser handles TXT file parsing
type TXTParser struct{}

// NewTXTParser creates a new TXT parser
func NewTXTParser() *TXTParser {
	return &TXTParser{}
}

// Parse extracts text from TXT file
func (p *TXTParser) Parse(reader io.Reader) (string, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// SupportedType returns the file type this parser supports
func (p *TXTParser) SupportedType() string {
	return "txt"
}
