package parser

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTXTParser_SupportedType(t *testing.T) {
	parser := NewTXTParser()
	assert.Equal(t, "txt", parser.SupportedType())
}

func TestTXTParser_Parse_EmptyReader(t *testing.T) {
	parser := NewTXTParser()

	// Empty reader should return empty string without error
	emptyReader := bytes.NewReader([]byte{})
	text, err := parser.Parse(emptyReader)

	require.NoError(t, err)
	assert.Empty(t, text)
}

func TestTXTParser_Parse_SimpleText(t *testing.T) {
	parser := NewTXTParser()

	content := "Hello, World!"
	reader := bytes.NewReader([]byte(content))

	text, err := parser.Parse(reader)

	require.NoError(t, err)
	assert.Equal(t, content, text)
}

func TestTXTParser_Parse_MultilineText(t *testing.T) {
	parser := NewTXTParser()

	content := "Line 1\nLine 2\nLine 3"
	reader := bytes.NewReader([]byte(content))

	text, err := parser.Parse(reader)

	require.NoError(t, err)
	assert.Equal(t, content, text)
	assert.Contains(t, text, "\n")
}

func TestTXTParser_Parse_UnicodeContent(t *testing.T) {
	parser := NewTXTParser()

	// Test various unicode characters
	content := "Привет мир! 你好世界! مرحبا بالعالم! 🎉🚀"
	reader := bytes.NewReader([]byte(content))

	text, err := parser.Parse(reader)

	require.NoError(t, err)
	assert.Equal(t, content, text)
}

func TestTXTParser_Parse_SpecialChars(t *testing.T) {
	parser := NewTXTParser()

	content := "Special chars: <>&\"' \t\n\r"
	reader := bytes.NewReader([]byte(content))

	text, err := parser.Parse(reader)

	require.NoError(t, err)
	assert.Equal(t, content, text)
}

func TestTXTParser_Parse_Whitespace(t *testing.T) {
	parser := NewTXTParser()

	tests := []struct {
		name    string
		content string
	}{
		{"spaces only", "   "},
		{"tabs only", "\t\t\t"},
		{"newlines only", "\n\n\n"},
		{"mixed whitespace", "  \t\n  \t\n  "},
		{"text with leading whitespace", "   Hello"},
		{"text with trailing whitespace", "Hello   "},
		{"text with both", "   Hello   "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader([]byte(tt.content))
			text, err := parser.Parse(reader)

			require.NoError(t, err)
			assert.Equal(t, tt.content, text)
		})
	}
}

func TestTXTParser_Parse_LargeContent(t *testing.T) {
	parser := NewTXTParser()

	// Generate large content (1MB)
	var builder strings.Builder
	line := "This is a test line for large content parsing. "
	for builder.Len() < 1024*1024 {
		builder.WriteString(line)
	}
	content := builder.String()

	reader := bytes.NewReader([]byte(content))
	text, err := parser.Parse(reader)

	require.NoError(t, err)
	assert.Equal(t, content, text)
	assert.True(t, len(text) >= 1024*1024)
}

func TestTXTParser_Parse_BinaryContent(t *testing.T) {
	parser := NewTXTParser()

	// Binary content with null bytes
	content := []byte{0x00, 0x01, 0x02, 'H', 'e', 'l', 'l', 'o', 0x00}
	reader := bytes.NewReader(content)

	text, err := parser.Parse(reader)

	require.NoError(t, err)
	assert.Equal(t, string(content), text)
}

func TestTXTParser_Parse_WindowsLineEndings(t *testing.T) {
	parser := NewTXTParser()

	content := "Line 1\r\nLine 2\r\nLine 3"
	reader := bytes.NewReader([]byte(content))

	text, err := parser.Parse(reader)

	require.NoError(t, err)
	assert.Equal(t, content, text)
	assert.Contains(t, text, "\r\n")
}

func TestTXTParser_Parse_MixedLineEndings(t *testing.T) {
	parser := NewTXTParser()

	content := "Unix line\nWindows line\r\nOld Mac line\rAnother line"
	reader := bytes.NewReader([]byte(content))

	text, err := parser.Parse(reader)

	require.NoError(t, err)
	assert.Equal(t, content, text)
}

// errorReader is a mock reader that always returns an error
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func TestTXTParser_Parse_ReaderError(t *testing.T) {
	parser := NewTXTParser()

	// Test with a reader that returns an error
	reader := &errorReader{}
	_, err := parser.Parse(reader)

	require.Error(t, err)
	assert.Equal(t, io.ErrUnexpectedEOF, err)
}

func TestTXTParser_Parse_UTF8BOM(t *testing.T) {
	parser := NewTXTParser()

	// UTF-8 BOM followed by content
	bom := []byte{0xEF, 0xBB, 0xBF}
	content := append(bom, []byte("Hello with BOM")...)
	reader := bytes.NewReader(content)

	text, err := parser.Parse(reader)

	require.NoError(t, err)
	// Parser preserves BOM as-is
	assert.Equal(t, string(content), text)
}

func TestTXTParser_Parse_UTF16Content(t *testing.T) {
	parser := NewTXTParser()

	// UTF-16 LE BOM + "Hi" in UTF-16 LE
	content := []byte{0xFF, 0xFE, 'H', 0x00, 'i', 0x00}
	reader := bytes.NewReader(content)

	text, err := parser.Parse(reader)

	require.NoError(t, err)
	// Parser reads bytes as-is (no encoding conversion)
	assert.Equal(t, string(content), text)
}

func TestTXTParser_Parse_JSONContent(t *testing.T) {
	parser := NewTXTParser()

	content := `{"key": "value", "number": 123, "array": [1, 2, 3]}`
	reader := bytes.NewReader([]byte(content))

	text, err := parser.Parse(reader)

	require.NoError(t, err)
	assert.Equal(t, content, text)
}

func TestTXTParser_Parse_XMLContent(t *testing.T) {
	parser := NewTXTParser()

	content := `<?xml version="1.0"?><root><item>value</item></root>`
	reader := bytes.NewReader([]byte(content))

	text, err := parser.Parse(reader)

	require.NoError(t, err)
	assert.Equal(t, content, text)
}

func TestTXTParser_Parse_CodeContent(t *testing.T) {
	parser := NewTXTParser()

	content := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`
	reader := bytes.NewReader([]byte(content))

	text, err := parser.Parse(reader)

	require.NoError(t, err)
	assert.Equal(t, content, text)
}

func TestTXTParser_Parse_MarkdownContent(t *testing.T) {
	parser := NewTXTParser()

	content := `# Header

## Subheader

- Item 1
- Item 2

**Bold** and *italic* text.

` + "```go\ncode block\n```"

	reader := bytes.NewReader([]byte(content))

	text, err := parser.Parse(reader)

	require.NoError(t, err)
	assert.Equal(t, content, text)
}

// Benchmark tests
func BenchmarkTXTParser_Parse_Small(b *testing.B) {
	parser := NewTXTParser()
	content := []byte("Small content for benchmark testing.")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(bytes.NewReader(content))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTXTParser_Parse_Medium(b *testing.B) {
	parser := NewTXTParser()

	// 100KB content
	var builder strings.Builder
	for builder.Len() < 100*1024 {
		builder.WriteString("Medium content line for benchmark. ")
	}
	content := []byte(builder.String())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(bytes.NewReader(content))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTXTParser_Parse_Large(b *testing.B) {
	parser := NewTXTParser()

	// 1MB content
	var builder strings.Builder
	for builder.Len() < 1024*1024 {
		builder.WriteString("Large content line for benchmark testing. ")
	}
	content := []byte(builder.String())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(bytes.NewReader(content))
		if err != nil {
			b.Fatal(err)
		}
	}
}
