package parser

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	ledongthucpdf "github.com/ledongthuc/pdf"
	rscpdf "rsc.io/pdf"
)

// PDFParser handles PDF file parsing
type PDFParser struct{}

// NewPDFParser creates a new PDF parser
func NewPDFParser() *PDFParser {
	return &PDFParser{}
}

// Parse extracts text from PDF file using multiple parsing strategies
func (p *PDFParser) Parse(reader io.Reader) (string, error) {
	// Read all bytes from the reader into a buffer
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read PDF data: %w", err)
	}

	// Try rsc.io/pdf first (simple but effective for many PDFs)
	text, err := p.parseWithRscPDF(data)
	if err == nil && len(strings.TrimSpace(text)) > 0 {
		return p.cleanText(text), nil
	}

	// Fallback to ledongthuc/pdf (better support for complex structures)
	text, err = p.parseWithLedongthuc(data)
	if err == nil && len(strings.TrimSpace(text)) > 0 {
		return p.cleanText(text), nil
	}

	// If both methods failed, return informative metadata
	numPages := p.getPageCount(data)
	metadata := fmt.Sprintf("[PDF Document Metadata]\n"+
		"Pages: %d\n"+
		"Note: Text extraction failed with multiple parsers. This PDF may contain:\n"+
		"- Scanned images without OCR text layer\n"+
		"- Complex font encoding not supported by available libraries\n"+
		"- Embedded images or graphics instead of text\n\n"+
		"Suggestion: Please try uploading a text-based PDF or use OCR software to convert scanned PDFs.",
		numPages)
	return metadata, nil
}

// parseWithRscPDF uses rsc.io/pdf library for text extraction
func (p *PDFParser) parseWithRscPDF(data []byte) (string, error) {
	// Create a bytes reader
	bytesReader := bytes.NewReader(data)

	// Open the PDF document
	pdfReader, err := rscpdf.NewReader(bytesReader, int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("rsc.io/pdf: failed to open PDF: %w", err)
	}

	numPages := pdfReader.NumPage()
	var builder strings.Builder

	// Extract text from each page
	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page := pdfReader.Page(pageNum)
		if page.V.IsNull() {
			continue
		}

		// Get page content
		content := page.Content()

		// Check if there's any text content
		if len(content.Text) == 0 {
			continue
		}

		// Add page separator
		if pageNum > 1 {
			builder.WriteString("\n\n--- Page ")
			builder.WriteString(fmt.Sprintf("%d", pageNum))
			builder.WriteString(" ---\n\n")
		}

		// Extract text from all Text objects on the page
		for _, text := range content.Text {
			builder.WriteString(text.S)
			builder.WriteString(" ") // Add space between text fragments
		}
	}

	return builder.String(), nil
}

// parseWithLedongthuc uses ledongthuc/pdf library for text extraction
func (p *PDFParser) parseWithLedongthuc(data []byte) (string, error) {
	// Create a bytes reader for the PDF library
	bytesReader := bytes.NewReader(data)

	// Open the PDF document
	pdfReader, err := ledongthucpdf.NewReader(bytesReader, int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("ledongthuc: failed to open PDF: %w", err)
	}

	// Extract text from all pages
	var builder strings.Builder
	numPages := pdfReader.NumPage()

	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page := pdfReader.Page(pageNum)
		if page.V.IsNull() {
			continue
		}

		// Extract text content from the page
		text, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}

		// Add page separator for readability
		if pageNum > 1 && len(text) > 0 {
			builder.WriteString("\n\n--- Page ")
			builder.WriteString(fmt.Sprintf("%d", pageNum))
			builder.WriteString(" ---\n\n")
		}

		builder.WriteString(text)
	}

	return builder.String(), nil
}

// getPageCount tries to get page count from PDF
func (p *PDFParser) getPageCount(data []byte) int {
	// Try with rsc.io/pdf first
	bytesReader := bytes.NewReader(data)
	rscReader, err := rscpdf.NewReader(bytesReader, int64(len(data)))
	if err == nil {
		return rscReader.NumPage()
	}

	// Fallback to ledongthuc
	bytesReader = bytes.NewReader(data)
	ledongthucReader, err := ledongthucpdf.NewReader(bytesReader, int64(len(data)))
	if err == nil {
		return ledongthucReader.NumPage()
	}

	return 0

}

// cleanText normalizes extracted text by removing extra whitespace and empty lines
func (p *PDFParser) cleanText(text string) string {
	// Replace multiple spaces/tabs with single space (but preserve newlines)
	re := regexp.MustCompile(`[^\S\n]+`)
	text = re.ReplaceAllString(text, " ")

	// Remove spaces at the beginning and end of each line
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	text = strings.Join(lines, "\n")

	// Replace 3+ consecutive newlines with double newline
	re = regexp.MustCompile(`\n{3,}`)
	text = re.ReplaceAllString(text, "\n\n")

	return strings.TrimSpace(text)
}

// SupportedType returns the file type this parser supports
func (p *PDFParser) SupportedType() string {
	return "pdf"
}
