// Package validate provides code validation utilities.
// syntax.go: Basic syntax validation.
// DACE: Reusable validation functions.
// P2 FIX: Skip braces inside strings and comments.
package validate

import "strings"

// GoSyntax validates basic Go syntax (braces, parens).
// P2 FIX: Now correctly skips braces inside strings and comments.
func GoSyntax(content string) string {
	braceCount := 0
	parenCount := 0
	inString := false
	inRawString := false
	inLineComment := false
	inBlockComment := false
	prevChar := rune(0)

	runes := []rune(content)
	for i, ch := range runes {
		// Handle newlines - end line comments
		if ch == '\n' {
			inLineComment = false
			prevChar = ch
			continue
		}

		// Skip if in line comment
		if inLineComment {
			prevChar = ch
			continue
		}

		// Handle block comment end
		if inBlockComment {
			if prevChar == '*' && ch == '/' {
				inBlockComment = false
			}
			prevChar = ch
			continue
		}

		// Check for comment start (only if not in string)
		if !inString && !inRawString && ch == '/' && i+1 < len(runes) {
			if runes[i+1] == '/' {
				inLineComment = true
				prevChar = ch
				continue
			}
			if runes[i+1] == '*' {
				inBlockComment = true
				prevChar = ch
				continue
			}
		}

		// Handle raw string (backtick)
		if ch == '`' && !inString {
			inRawString = !inRawString
			prevChar = ch
			continue
		}

		// Handle regular string
		if ch == '"' && !inRawString && prevChar != '\\' {
			inString = !inString
			prevChar = ch
			continue
		}

		// Skip braces inside strings
		if inString || inRawString {
			prevChar = ch
			continue
		}

		// Count braces only outside strings and comments
		switch ch {
		case '{':
			braceCount++
		case '}':
			braceCount--
		case '(':
			parenCount++
		case ')':
			parenCount--
		}
		prevChar = ch
	}

	if braceCount != 0 {
		return "unbalanced braces"
	}
	if parenCount != 0 {
		return "unbalanced parentheses"
	}
	return ""
}

// JSONSyntax validates basic JSON syntax.
// P2 FIX: Now correctly skips braces inside strings.
func JSONSyntax(content string) string {
	content = strings.TrimSpace(content)
	if len(content) == 0 {
		return "empty JSON"
	}

	if content[0] != '{' && content[0] != '[' {
		return "must start with { or ["
	}

	braceCount := 0
	bracketCount := 0
	inString := false
	prevChar := rune(0)

	for _, ch := range content {
		// Handle string boundaries
		if ch == '"' && prevChar != '\\' {
			inString = !inString
			prevChar = ch
			continue
		}

		// Skip braces inside strings
		if inString {
			prevChar = ch
			continue
		}

		switch ch {
		case '{':
			braceCount++
		case '}':
			braceCount--
		case '[':
			bracketCount++
		case ']':
			bracketCount--
		}
		prevChar = ch
	}

	if braceCount != 0 {
		return "unbalanced braces"
	}
	if bracketCount != 0 {
		return "unbalanced brackets"
	}
	return ""
}
