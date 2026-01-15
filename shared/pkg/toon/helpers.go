// Package toon provides TOON parsing utilities.
// helpers.go: Parser helper functions.
// DACE: Single responsibility - helper functions only.
package toon

// isBlockHeader checks if a key looks like a TOON block header.
// Block headers are typically UPPERCASE or Title_Case with underscores.
// Examples: "RESEARCH", "META", "DECISION", "CEO_DECISIONS", "FACT"
func isBlockHeader(key string) bool {
	if len(key) == 0 {
		return false
	}
	// Must start with uppercase letter
	if key[0] < 'A' || key[0] > 'Z' {
		return false
	}
	// Count uppercase letters - block headers have high uppercase ratio
	upperCount := 0
	for _, ch := range key {
		if ch >= 'A' && ch <= 'Z' {
			upperCount++
		}
	}
	// At least 50% uppercase or all uppercase
	return upperCount*2 >= len(key) || upperCount == len(key)
}
