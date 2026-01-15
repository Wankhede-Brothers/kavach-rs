// Package toon provides TOON parsing.
// parser_test.go: Tests for TOON parser.
package toon

import (
	"strings"
	"testing"
)

func TestParser_Parse_BracketFormat(t *testing.T) {
	input := `
[DECISION]
id: D001
status: active
desc: Test decision

[META]
version: 1.0
`
	p := NewParser()
	doc, err := p.ParseString(input)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(doc.Blocks) != 2 {
		t.Errorf("expected 2 blocks, got %d", len(doc.Blocks))
	}

	decision, ok := doc.Blocks["DECISION"]
	if !ok {
		t.Fatal("DECISION block not found")
	}

	if decision.Fields["id"] != "D001" {
		t.Errorf("expected id=D001, got %s", decision.Fields["id"])
	}

	if decision.Fields["status"] != "active" {
		t.Errorf("expected status=active, got %s", decision.Fields["status"])
	}
}

func TestParser_Parse_SP30Format(t *testing.T) {
	input := `
RESEARCH:active
  topic: kubernetes
  date: 2026-01-17

META:system
  version: 3.0
`
	p := NewParser()
	doc, err := p.ParseString(input)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	research, ok := doc.Blocks["RESEARCH"]
	if !ok {
		t.Fatal("RESEARCH block not found")
	}

	if research.Fields["_value"] != "active" {
		t.Errorf("expected _value=active, got %s", research.Fields["_value"])
	}

	if research.Fields["topic"] != "kubernetes" {
		t.Errorf("expected topic=kubernetes, got %s", research.Fields["topic"])
	}
}

func TestParser_Parse_Arrays(t *testing.T) {
	input := `
[PATTERNS]
keywords[]:
  - debug
  - error
  - fix
`
	p := NewParser()
	doc, err := p.ParseString(input)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	patterns, ok := doc.Blocks["PATTERNS"]
	if !ok {
		t.Fatal("PATTERNS block not found")
	}

	keywords := patterns.Arrays["keywords"]
	if len(keywords) != 3 {
		t.Errorf("expected 3 keywords, got %d", len(keywords))
	}

	if keywords[0] != "debug" || keywords[1] != "error" || keywords[2] != "fix" {
		t.Errorf("unexpected keywords: %v", keywords)
	}
}

func TestParser_Parse_TypedArrays(t *testing.T) {
	input := `
[CONFIG]
tools[8]{string}: Read,Write,Edit
items[3]: first
`
	p := NewParser()
	doc, err := p.ParseString(input)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	config, ok := doc.Blocks["CONFIG"]
	if !ok {
		t.Fatal("CONFIG block not found")
	}

	if config.ArraySizes["tools"] != 8 {
		t.Errorf("expected tools size 8, got %d", config.ArraySizes["tools"])
	}

	if config.ArrayTypes["tools"] != "string" {
		t.Errorf("expected tools type string, got %s", config.ArrayTypes["tools"])
	}

	if config.ArraySizes["items"] != 3 {
		t.Errorf("expected items size 3, got %d", config.ArraySizes["items"])
	}
}

func TestParser_Parse_Comments(t *testing.T) {
	input := `
# This is a comment
[BLOCK]
# Another comment
key: value
`
	p := NewParser()
	doc, err := p.ParseString(input)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(doc.Blocks) != 1 {
		t.Errorf("expected 1 block, got %d", len(doc.Blocks))
	}

	block := doc.Blocks["BLOCK"]
	if block.Fields["key"] != "value" {
		t.Errorf("expected key=value, got %s", block.Fields["key"])
	}
}

func TestParser_Parse_EmptyInput(t *testing.T) {
	p := NewParser()
	doc, err := p.ParseString("")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(doc.Blocks) != 0 {
		t.Errorf("expected 0 blocks for empty input, got %d", len(doc.Blocks))
	}
}

func TestParser_Parse_OnlyComments(t *testing.T) {
	input := `
# Comment 1
# Comment 2
`
	p := NewParser()
	doc, err := p.ParseString(input)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(doc.Blocks) != 0 {
		t.Errorf("expected 0 blocks, got %d", len(doc.Blocks))
	}
}

func TestParseArrayKey(t *testing.T) {
	tests := []struct {
		key      string
		wantName string
		wantSize int
		wantType string
	}{
		{"keywords[]", "keywords", 0, ""},
		{"keywords[3]", "keywords", 3, ""},
		{"keywords[3]{string}", "keywords", 3, "string"},
		{"keywords{string}", "keywords", 0, "string"},
		{"tools[8]{string}", "tools", 8, "string"},
		{"items", "items", 0, ""},
		{"arr[10]{int}", "arr", 10, "int"},
	}

	for _, tt := range tests {
		name, size, typ := parseArrayKey(tt.key)
		if name != tt.wantName {
			t.Errorf("parseArrayKey(%q).name = %s, want %s", tt.key, name, tt.wantName)
		}
		if size != tt.wantSize {
			t.Errorf("parseArrayKey(%q).size = %d, want %d", tt.key, size, tt.wantSize)
		}
		if typ != tt.wantType {
			t.Errorf("parseArrayKey(%q).typ = %s, want %s", tt.key, typ, tt.wantType)
		}
	}
}

func TestIsBlockHeader(t *testing.T) {
	tests := []struct {
		key  string
		want bool
	}{
		{"RESEARCH", true},
		{"META", true},
		{"DECISION", true},
		{"CEO_DECISIONS", true},
		{"Fact", false}, // Only 1/4 uppercase
		{"research", false},
		{"key", false},
		{"", false},
		{"A", true},
		{"ABC", true},
		{"ABc", true},  // 2/3 uppercase
		{"Abc", false}, // 1/3 uppercase
	}

	for _, tt := range tests {
		got := isBlockHeader(tt.key)
		if got != tt.want {
			t.Errorf("isBlockHeader(%q) = %v, want %v", tt.key, got, tt.want)
		}
	}
}

func TestParser_Parse_MixedFormats(t *testing.T) {
	input := `
[BRACKET_BLOCK]
key1: value1

SP30_BLOCK:header_value
  key2: value2
`
	p := NewParser()
	doc, err := p.ParseString(input)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(doc.Blocks) != 2 {
		t.Errorf("expected 2 blocks, got %d", len(doc.Blocks))
	}

	bracketBlock := doc.Blocks["BRACKET_BLOCK"]
	if bracketBlock.Fields["key1"] != "value1" {
		t.Errorf("expected key1=value1, got %s", bracketBlock.Fields["key1"])
	}

	sp30Block := doc.Blocks["SP30_BLOCK"]
	if sp30Block.Fields["_value"] != "header_value" {
		t.Errorf("expected _value=header_value, got %s", sp30Block.Fields["_value"])
	}
}

func TestParser_Parse_ValuesWithColons(t *testing.T) {
	input := `
[CONFIG]
url: https://example.com:8080/path
time: 12:30:45
`
	p := NewParser()
	doc, err := p.ParseString(input)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	config := doc.Blocks["CONFIG"]
	if config.Fields["url"] != "https://example.com:8080/path" {
		t.Errorf("expected url with colons, got %s", config.Fields["url"])
	}

	if config.Fields["time"] != "12:30:45" {
		t.Errorf("expected time with colons, got %s", config.Fields["time"])
	}
}

func TestParser_Parse_LargeDocument(t *testing.T) {
	// Build a large document
	var builder strings.Builder
	for i := 0; i < 100; i++ {
		builder.WriteString("[BLOCK")
		builder.WriteString(string(rune('A' + i%26)))
		builder.WriteString("]\n")
		builder.WriteString("key: value")
		builder.WriteString(string(rune('0' + i%10)))
		builder.WriteString("\n")
	}

	p := NewParser()
	doc, err := p.ParseString(builder.String())
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Due to duplicate names, we may have fewer blocks
	if len(doc.Blocks) < 26 {
		t.Errorf("expected at least 26 unique blocks, got %d", len(doc.Blocks))
	}
}
