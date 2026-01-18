// Package toon provides TOON (Token-Oriented Object Notation) parsing.
// types.go: Core data types for TOON documents.
// DACE: Single responsibility - type definitions only.
package toon

// Block represents a TOON block with key-value pairs.
// P2 FIX #6: Added ArrayTypes for SP/1.0 typed arrays support.
type Block struct {
	Name       string
	Fields     map[string]string
	Arrays     map[string][]string
	ArrayTypes map[string]string // P2 FIX #6: Stores type annotation like "string" from keywords[3]{string}
	ArraySizes map[string]int    // P2 FIX #6: Stores size annotation like 3 from keywords[3]{string}
	Nested     map[string]*Block
}

// NewBlock creates an empty TOON block.
func NewBlock(name string) *Block {
	return &Block{
		Name:       name,
		Fields:     make(map[string]string),
		Arrays:     make(map[string][]string),
		ArrayTypes: make(map[string]string),
		ArraySizes: make(map[string]int),
		Nested:     make(map[string]*Block),
	}
}

// Get returns a field value or empty string.
func (b *Block) Get(key string) string {
	if b.Fields == nil {
		return ""
	}
	return b.Fields[key]
}

// GetArray returns an array value or empty slice.
func (b *Block) GetArray(key string) []string {
	if b.Arrays == nil {
		return nil
	}
	return b.Arrays[key]
}

// GetNested returns a nested block or nil.
func (b *Block) GetNested(key string) *Block {
	if b.Nested == nil {
		return nil
	}
	return b.Nested[key]
}

// GetArrayType returns the type annotation for an array or empty string.
// P2 FIX #6: For keywords[3]{string}, returns "string".
func (b *Block) GetArrayType(key string) string {
	if b.ArrayTypes == nil {
		return ""
	}
	return b.ArrayTypes[key]
}

// GetArraySize returns the size annotation for an array or 0.
// P2 FIX #6: For keywords[3]{string}, returns 3.
func (b *Block) GetArraySize(key string) int {
	if b.ArraySizes == nil {
		return 0
	}
	return b.ArraySizes[key]
}

// Document represents a parsed TOON document.
type Document struct {
	Blocks map[string]*Block
}

// NewDocument creates an empty TOON document.
func NewDocument() *Document {
	return &Document{
		Blocks: make(map[string]*Block),
	}
}

// Get returns a block by name.
func (d *Document) Get(name string) *Block {
	return d.Blocks[name]
}
