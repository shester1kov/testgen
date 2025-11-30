package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "plain text",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "script tag",
			input:    "Hello<script>alert('XSS')</script>World",
			expected: "HelloWorld",
		},
		{
			name:     "script tag uppercase",
			input:    "Test<SCRIPT>alert('XSS')</SCRIPT>End",
			expected: "TestEnd",
		},
		{
			name:     "onclick event",
			input:    "Click<div onclick='alert()'>here</div>",
			expected: "Click<div>here</div>",
		},
		{
			name:     "onerror event",
			input:    "Image<img onerror='alert()' src='x'>",
			expected: "Image<img src='x'>",
		},
		{
			name:     "javascript protocol",
			input:    "Link<a href='javascript:alert()'>click</a>",
			expected: "Link<a href='alert()'>click</a>",
		},
		{
			name:     "data protocol",
			input:    "Link<a href='data:text/html,<script>alert()</script>'>click</a>",
			expected: "Link<a href=','>click</a>", // Both data:text/html and script tag removed
		},
		{
			name:     "multiple XSS vectors",
			input:    "<script>alert()</script>Test<div onclick='bad()'>text</div>",
			expected: "Test<div>text</div>",
		},
		{
			name:     "whitespace trimming",
			input:    "  Test  ",
			expected: "Test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeInput(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSanitizeHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "plain text",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "HTML entities",
			input:    "<script>alert('XSS')</script>",
			expected: "&lt;script&gt;alert(&#39;XSS&#39;)&lt;/script&gt;",
		},
		{
			name:     "ampersand",
			input:    "Tom & Jerry",
			expected: "Tom &amp; Jerry",
		},
		{
			name:     "quotes",
			input:    `He said "Hello"`,
			expected: "He said &#34;Hello&#34;",
		},
		{
			name:     "mixed content",
			input:    `<div onclick="alert('XSS')">Test & "quotes"</div>`,
			expected: `&lt;div onclick=&#34;alert(&#39;XSS&#39;)&#34;&gt;Test &amp; &#34;quotes&#34;&lt;/div&gt;`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeHTML(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSanitizeMultiline(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single line",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "multiline plain text",
			input:    "Line 1\nLine 2\nLine 3",
			expected: "Line 1\nLine 2\nLine 3",
		},
		{
			name:     "multiline with XSS",
			input:    "Line 1<script>alert()</script>\nLine 2<div onclick='bad()'>test</div>\nLine 3",
			expected: "Line 1\nLine 2<div>test</div>\nLine 3",
		},
		{
			name:     "multiline with whitespace",
			input:    "  Line 1  \n  Line 2  ",
			expected: "Line 1\nLine 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeMultiline(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
