package validate

import "testing"

func TestGoSyntax(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr string
	}{
		{"valid simple", "func main() { }", ""},
		{"valid with string braces", `fmt.Println("Hello {world}")`, ""},
		{"valid with comment braces", "// {comment}\nfunc f() {}", ""},
		{"valid with block comment", "/* {block} */\nfunc f() {}", ""},
		{"valid with raw string", "x := `{raw}`", ""},
		{"unbalanced open", "func main() {", "unbalanced braces"},
		{"unbalanced close", "func main() }", "unbalanced braces"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GoSyntax(tt.content)
			if got != tt.wantErr {
				t.Errorf("GoSyntax() = %q, want %q", got, tt.wantErr)
			}
		})
	}
}

func TestJSONSyntax(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr string
	}{
		{"valid object", `{"key": "value"}`, ""},
		{"valid with string braces", `{"key": "{value}"}`, ""},
		{"valid array", `[1, 2, 3]`, ""},
		{"empty", "", "empty JSON"},
		{"not json", "hello", "must start with { or ["},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := JSONSyntax(tt.content)
			if got != tt.wantErr {
				t.Errorf("JSONSyntax() = %q, want %q", got, tt.wantErr)
			}
		})
	}
}
