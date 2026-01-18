// Package toon provides TOON (Token-Oriented Object Notation) parsing.
// parser.go: Document parsing logic.
// DACE: Single responsibility - parsing only.
// TOON achieves ~40% token savings vs JSON through compact syntax.
package toon

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

// Parser provides TOON document parsing.
type Parser struct{}

// NewParser creates a new TOON parser.
func NewParser() *Parser {
	return &Parser{}
}

// Parse reads and parses a TOON document from a reader.
// Supports two TOON formats:
//  1. Bracket format: [BLOCK_NAME] followed by key: value
//  2. SP/1.0 format: BLOCK_NAME:value followed by indented key: value
func (p *Parser) Parse(r io.Reader) (*Document, error) {
	doc := NewDocument()
	scanner := bufio.NewScanner(r)

	var currentBlock *Block
	var currentArray string

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Skip empty lines and comments
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Format 1: Block header [BLOCK_NAME]
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			blockName := trimmed[1 : len(trimmed)-1]
			currentBlock = NewBlock(blockName)
			doc.Blocks[blockName] = currentBlock
			currentArray = ""
			continue
		}

		// Format 2: SP/1.0 block header BLOCK_NAME:value
		if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
			if idx := strings.Index(trimmed, ":"); idx > 0 {
				key := trimmed[:idx]
				if isBlockHeader(key) {
					blockName := key
					value := strings.TrimSpace(trimmed[idx+1:])
					currentBlock = NewBlock(blockName)
					doc.Blocks[blockName] = currentBlock
					currentArray = ""
					if value != "" {
						currentBlock.Fields["_value"] = value
					}
					continue
				}
			}
		}

		if currentBlock == nil {
			continue
		}

		// Array continuation
		if currentArray != "" && (strings.HasPrefix(line, "  ") || strings.HasPrefix(line, "\t")) {
			item := strings.TrimSpace(trimmed)
			if strings.HasPrefix(item, "- ") {
				item = strings.TrimPrefix(item, "- ")
			}
			currentBlock.Arrays[currentArray] = append(currentBlock.Arrays[currentArray], item)
			continue
		}

		// Indented key-value pair (SP/1.0 format)
		if strings.HasPrefix(line, "  ") || strings.HasPrefix(line, "\t") {
			if idx := strings.Index(trimmed, ":"); idx > 0 {
				key := strings.TrimSpace(trimmed[:idx])
				value := strings.TrimSpace(trimmed[idx+1:])
				if strings.Contains(key, "[") {
					// P2 FIX #6: Parse typed arrays like keywords[3]{string}
					arrayName, size, typ := parseArrayKey(key)
					currentArray = arrayName
					currentBlock.Arrays[arrayName] = []string{}
					if size > 0 {
						currentBlock.ArraySizes[arrayName] = size
					}
					if typ != "" {
						currentBlock.ArrayTypes[arrayName] = typ
					}
					if value != "" {
						currentBlock.Arrays[arrayName] = append(currentBlock.Arrays[arrayName], value)
					}
					continue
				}
				currentBlock.Fields[key] = value
				currentArray = ""
				continue
			}
		}

		// Non-indented key-value pair (bracket format)
		if idx := strings.Index(trimmed, ":"); idx > 0 {
			key := strings.TrimSpace(trimmed[:idx])
			value := strings.TrimSpace(trimmed[idx+1:])
			// P2 FIX #6: Check for typed array syntax keywords[3]{string} or keywords[]
			if strings.Contains(key, "[") {
				arrayName, size, typ := parseArrayKey(key)
				currentArray = arrayName
				currentBlock.Arrays[arrayName] = []string{}
				if size > 0 {
					currentBlock.ArraySizes[arrayName] = size
				}
				if typ != "" {
					currentBlock.ArrayTypes[arrayName] = typ
				}
				if value != "" {
					currentBlock.Arrays[arrayName] = append(currentBlock.Arrays[arrayName], value)
				}
				continue
			}
			currentBlock.Fields[key] = value
			currentArray = ""
		}
	}

	return doc, scanner.Err()
}

// parseArrayKey parses array key syntax like:
//   - keywords[] -> ("keywords", 0, "")
//   - keywords[3] -> ("keywords", 3, "")
//   - keywords[3]{string} -> ("keywords", 3, "string")
//   - keywords{string} -> ("keywords", 0, "string")
//
// P2 FIX #6: Adds support for SP/1.0 typed arrays.
// P2 FIX: Handle type-only syntax (keywords{string}) without brackets.
func parseArrayKey(key string) (name string, size int, typ string) {
	// Find bracket position
	bracketIdx := strings.Index(key, "[")
	braceIdx := strings.Index(key, "{")

	// P2 FIX: Handle type-only syntax like keywords{string}
	if bracketIdx == -1 && braceIdx == -1 {
		return key, 0, ""
	}

	// Determine name end position (first of [ or {)
	nameEnd := len(key)
	if bracketIdx != -1 {
		nameEnd = bracketIdx
	}
	if braceIdx != -1 && braceIdx < nameEnd {
		nameEnd = braceIdx
	}
	name = key[:nameEnd]

	// Extract size from [N]
	if bracketIdx != -1 {
		closeBracket := strings.Index(key, "]")
		if closeBracket > bracketIdx+1 {
			sizeStr := key[bracketIdx+1 : closeBracket]
			if s, err := strconv.Atoi(sizeStr); err == nil {
				size = s
			}
		}
	}

	// Extract type from {type}
	if braceIdx != -1 {
		closeBrace := strings.Index(key, "}")
		if closeBrace > braceIdx+1 {
			typ = key[braceIdx+1 : closeBrace]
		}
	}

	return name, size, typ
}

// ParseString parses a TOON document from a string.
func (p *Parser) ParseString(s string) (*Document, error) {
	return p.Parse(strings.NewReader(s))
}
