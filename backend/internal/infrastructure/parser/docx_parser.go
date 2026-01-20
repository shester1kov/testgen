package parser

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/nguyenthenguyen/docx"
)

// DOCXParser handles DOCX file parsing
type DOCXParser struct{}

// NewDOCXParser creates a new DOCX parser
func NewDOCXParser() *DOCXParser {
	return &DOCXParser{}
}

// Parse extracts text from DOCX file using nguyenthenguyen/docx library
func (p *DOCXParser) Parse(reader io.Reader) (string, error) {
	// Read all bytes from the reader
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read DOCX data: %w", err)
	}

	if len(data) == 0 {
		return "", fmt.Errorf("DOCX file is empty")
	}

	// Create a temporary file to store the DOCX data
	tmpFile, err := os.CreateTemp("", "docx-*.docx")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write the data to the temporary file
	if _, err := tmpFile.Write(data); err != nil {
		return "", fmt.Errorf("failed to write to temporary file: %w", err)
	}
	tmpFile.Close()

	// Open and parse the DOCX file
	doc, err := docx.ReadDocxFile(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to open DOCX document: %w", err)
	}
	defer doc.Close()

	// Extract text content (returns XML)
	docText := doc.Editable()
	if docText == nil {
		return "", fmt.Errorf("failed to get editable content from DOCX")
	}

	// Get the XML content
	xmlContent := docText.GetContent()

	// Extract plain text from XML
	text := p.extractTextFromXML(xmlContent)
	text = strings.TrimSpace(text)

	// Validate that we extracted some text
	if len(text) == 0 {
		return "", fmt.Errorf("no text content found in DOCX document")
	}

	return text, nil
}

// extractTextFromXML extracts plain text from Word XML content
func (p *DOCXParser) extractTextFromXML(xmlContent string) string {
	type TextElement struct {
		Text string `xml:",chardata"`
	}

	decoder := xml.NewDecoder(strings.NewReader(xmlContent))
	var builder strings.Builder
	var currentText strings.Builder

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}

		switch elem := token.(type) {
		case xml.StartElement:
			// Check if this is a text element (w:t)
			if elem.Name.Local == "t" {
				// Next token should be CharData with the text
			}
		case xml.CharData:
			// Collect character data
			text := strings.TrimSpace(string(elem))
			if len(text) > 0 {
				currentText.WriteString(text)
				currentText.WriteString(" ")
			}
		case xml.EndElement:
			// Add newline after paragraph (w:p)
			if elem.Name.Local == "p" && currentText.Len() > 0 {
				builder.WriteString(strings.TrimSpace(currentText.String()))
				builder.WriteString("\n")
				currentText.Reset()
			}
		}
	}

	// Add any remaining text
	if currentText.Len() > 0 {
		builder.WriteString(strings.TrimSpace(currentText.String()))
	}

	// Clean up extra whitespace
	text := builder.String()
	// Replace multiple spaces/tabs with single space (but preserve newlines)
	re := regexp.MustCompile(`[^\S\n]+`)
	text = re.ReplaceAllString(text, " ")
	// Replace multiple newlines with double newline (for paragraph separation)
	re = regexp.MustCompile(`\n{3,}`)
	text = re.ReplaceAllString(text, "\n\n")

	return strings.TrimSpace(text)
}

// ParseFromBytes is a helper method that parses DOCX from byte slice
func (p *DOCXParser) ParseFromBytes(data []byte) (string, error) {
	return p.Parse(bytes.NewReader(data))
}

// SupportedType returns the file type this parser supports
func (p *DOCXParser) SupportedType() string {
	return "docx"
}
