package security

import (
	"html"
	"regexp"
	"strings"
)

var (
	// scriptTagRegex matches script tags (case-insensitive)
	scriptTagRegex = regexp.MustCompile(`(?i)<script[\s\S]*?</script>`)

	// onEventRegex matches on* event handlers with their values
	onEventRegex = regexp.MustCompile(`(?i)\s+on\w+\s*=\s*['"][^'"]*['"]|\s+on\w+\s*=\s*[^\s>]+`)

	// javascriptProtocolRegex matches javascript: protocol
	javascriptProtocolRegex = regexp.MustCompile(`(?i)javascript:`)

	// dataProtocolRegex matches data: protocol (can be used for XSS)
	dataProtocolRegex = regexp.MustCompile(`(?i)data:text/html`)
)

// SanitizeInput sanitizes user input by removing XSS vectors
// This is a defense-in-depth measure - output encoding should still be used
func SanitizeInput(input string) string {
	if input == "" {
		return input
	}

	// Remove script tags
	input = scriptTagRegex.ReplaceAllString(input, "")

	// Remove event handlers (onclick, onerror, etc)
	input = onEventRegex.ReplaceAllString(input, "")

	// Remove javascript: protocol
	input = javascriptProtocolRegex.ReplaceAllString(input, "")

	// Remove dangerous data: protocols
	input = dataProtocolRegex.ReplaceAllString(input, "")

	// Trim whitespace
	input = strings.TrimSpace(input)

	return input
}

// SanitizeHTML escapes HTML entities to prevent XSS
// Use this for user-generated content that will be displayed as-is
func SanitizeHTML(input string) string {
	return html.EscapeString(input)
}

// SanitizeMultiline sanitizes multiline text (preserves newlines)
func SanitizeMultiline(input string) string {
	// First sanitize each line
	lines := strings.Split(input, "\n")
	for i, line := range lines {
		lines[i] = SanitizeInput(line)
	}
	return strings.Join(lines, "\n")
}
