package parser

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
)

// PPTXParser handles PPTX file parsing
type PPTXParser struct{}

// NewPPTXParser creates a new PPTX parser
func NewPPTXParser() *PPTXParser {
	return &PPTXParser{}
}

// Parse extracts text from PPTX file
// PPTX is a ZIP archive containing XML files for each slide
func (p *PPTXParser) Parse(reader io.Reader) (string, error) {
	// Read all bytes from the reader
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read PPTX data: %w", err)
	}

	if len(data) == 0 {
		return "", fmt.Errorf("PPTX file is empty")
	}

	// Open ZIP archive
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("failed to open PPTX as ZIP archive: %w", err)
	}

	// Find and sort slide files (ppt/slides/slide1.xml, slide2.xml, etc.)
	var slideFiles []*zip.File
	for _, file := range zipReader.File {
		if strings.HasPrefix(file.Name, "ppt/slides/slide") && strings.HasSuffix(file.Name, ".xml") {
			slideFiles = append(slideFiles, file)
		}
	}

	if len(slideFiles) == 0 {
		return "", fmt.Errorf("no slides found in PPTX document")
	}

	// Sort slides by number (slide1.xml, slide2.xml, ...)
	sort.Slice(slideFiles, func(i, j int) bool {
		return slideFiles[i].Name < slideFiles[j].Name
	})

	// Extract text from each slide
	var builder strings.Builder
	for i, slideFile := range slideFiles {
		text, err := p.extractTextFromSlide(slideFile)
		if err != nil {
			continue // Skip problematic slides
		}

		text = strings.TrimSpace(text)
		if len(text) > 0 {
			if builder.Len() > 0 {
				builder.WriteString("\n\n")
			}
			builder.WriteString(fmt.Sprintf("--- Slide %d ---\n", i+1))
			builder.WriteString(text)
		}
	}

	result := strings.TrimSpace(builder.String())
	if len(result) == 0 {
		return "", fmt.Errorf("no text content found in PPTX document")
	}

	return result, nil
}

// extractTextFromSlide extracts text from a single slide XML file
func (p *PPTXParser) extractTextFromSlide(file *zip.File) (string, error) {
	rc, err := file.Open()
	if err != nil {
		return "", err
	}
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil {
		return "", err
	}

	return p.extractTextFromXML(string(content)), nil
}

// extractTextFromXML extracts plain text from PowerPoint XML content
// PowerPoint uses similar structure to Word: <a:t> tags contain text, <a:p> is paragraph
func (p *PPTXParser) extractTextFromXML(xmlContent string) string {
	decoder := xml.NewDecoder(strings.NewReader(xmlContent))
	var builder strings.Builder
	var currentText strings.Builder
	inTextElement := false

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}

		switch elem := token.(type) {
		case xml.StartElement:
			// Text elements in PowerPoint XML: <a:t>
			if elem.Name.Local == "t" {
				inTextElement = true
			}
		case xml.CharData:
			if inTextElement {
				text := string(elem)
				if len(text) > 0 {
					currentText.WriteString(text)
				}
			}
		case xml.EndElement:
			if elem.Name.Local == "t" {
				inTextElement = false
			}
			// Add newline after paragraph (a:p) or text body paragraph
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
	// Replace multiple newlines with double newline
	re = regexp.MustCompile(`\n{3,}`)
	text = re.ReplaceAllString(text, "\n\n")

	return strings.TrimSpace(text)
}

// ParseFromBytes is a helper method that parses PPTX from byte slice
func (p *PPTXParser) ParseFromBytes(data []byte) (string, error) {
	return p.Parse(bytes.NewReader(data))
}

// SupportedType returns the file type this parser supports
func (p *PPTXParser) SupportedType() string {
	return "pptx"
}
