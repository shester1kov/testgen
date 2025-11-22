package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMDParser_Parse(t *testing.T) {
	parser := NewMDParser()

	t.Run("should parse markdown content successfully", func(t *testing.T) {
		markdown := `# Test Document

## Introduction

This is a **test** document with _markdown_ formatting.

### Features

- Feature 1
- Feature 2
- Feature 3

### Code Example

` + "```go\nfunc main() {\n    fmt.Println(\"Hello\")\n}\n```" + `

### Conclusion

This is the end of the document.
`
		reader := strings.NewReader(markdown)
		result, err := parser.Parse(reader)

		require.NoError(t, err)
		assert.Contains(t, result, "# Test Document")
		assert.Contains(t, result, "## Introduction")
		assert.Contains(t, result, "**test**")
		assert.Contains(t, result, "_markdown_")
		assert.Contains(t, result, "- Feature 1")
		assert.Contains(t, result, "func main()")
		assert.Contains(t, result, "Conclusion")
	})

	t.Run("should parse empty markdown", func(t *testing.T) {
		reader := strings.NewReader("")
		result, err := parser.Parse(reader)

		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("should parse markdown with unicode characters", func(t *testing.T) {
		markdown := "# Тестовый документ\n\n你好世界 Привет мир"
		reader := strings.NewReader(markdown)
		result, err := parser.Parse(reader)

		require.NoError(t, err)
		assert.Contains(t, result, "Тестовый документ")
		assert.Contains(t, result, "你好世界")
		assert.Contains(t, result, "Привет мир")
	})

	t.Run("should return correct supported type", func(t *testing.T) {
		assert.Equal(t, "md", parser.SupportedType())
	})
}
