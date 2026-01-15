// Package util provides utility functions.
// file_test.go: Tests for file utilities.
package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	existingFile := filepath.Join(tmpDir, "exists.txt")
	if err := os.WriteFile(existingFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		path string
		want bool
	}{
		{existingFile, true},
		{filepath.Join(tmpDir, "nonexistent.txt"), false},
		{tmpDir, true}, // directories also return true
	}

	for _, tt := range tests {
		got := FileExists(tt.path)
		if got != tt.want {
			t.Errorf("FileExists(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestDirExists(t *testing.T) {
	tmpDir := t.TempDir()
	existingFile := filepath.Join(tmpDir, "file.txt")
	if err := os.WriteFile(existingFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		path string
		want bool
	}{
		{tmpDir, true},
		{existingFile, false}, // file is not a directory
		{filepath.Join(tmpDir, "nonexistent"), false},
	}

	for _, tt := range tests {
		got := DirExists(tt.path)
		if got != tt.want {
			t.Errorf("DirExists(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestEnsureDir(t *testing.T) {
	tmpDir := t.TempDir()
	newDir := filepath.Join(tmpDir, "new", "nested", "dir")

	if err := EnsureDir(newDir); err != nil {
		t.Errorf("EnsureDir() error = %v", err)
	}

	if !DirExists(newDir) {
		t.Error("EnsureDir() did not create directory")
	}

	// Should be idempotent
	if err := EnsureDir(newDir); err != nil {
		t.Errorf("EnsureDir() second call error = %v", err)
	}
}

func TestEnsureParentDir(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "parent", "file.txt")

	if err := EnsureParentDir(filePath); err != nil {
		t.Errorf("EnsureParentDir() error = %v", err)
	}

	parentDir := filepath.Dir(filePath)
	if !DirExists(parentDir) {
		t.Error("EnsureParentDir() did not create parent directory")
	}
}

func TestReadWriteJSON(t *testing.T) {
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "test.json")

	type TestData struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	writeData := TestData{Name: "test", Count: 42}
	if err := WriteJSON(jsonFile, writeData); err != nil {
		t.Errorf("WriteJSON() error = %v", err)
	}

	var readData TestData
	if err := ReadJSON(jsonFile, &readData); err != nil {
		t.Errorf("ReadJSON() error = %v", err)
	}

	if readData.Name != writeData.Name || readData.Count != writeData.Count {
		t.Errorf("ReadJSON() = %+v, want %+v", readData, writeData)
	}

	// Test reading non-existent file
	if err := ReadJSON(filepath.Join(tmpDir, "nonexistent.json"), &readData); err == nil {
		t.Error("ReadJSON() expected error for non-existent file")
	}
}

func TestAppendFile(t *testing.T) {
	tmpDir := t.TempDir()
	appendFile := filepath.Join(tmpDir, "append.txt")

	if err := AppendFile(appendFile, []byte("line1\n")); err != nil {
		t.Errorf("AppendFile() first call error = %v", err)
	}

	if err := AppendFile(appendFile, []byte("line2\n")); err != nil {
		t.Errorf("AppendFile() second call error = %v", err)
	}

	content, err := os.ReadFile(appendFile)
	if err != nil {
		t.Fatal(err)
	}

	expected := "line1\nline2\n"
	if string(content) != expected {
		t.Errorf("AppendFile() content = %q, want %q", string(content), expected)
	}
}

func TestCopyFile(t *testing.T) {
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "source.txt")
	dstFile := filepath.Join(tmpDir, "dest", "copied.txt")

	content := []byte("test content")
	if err := os.WriteFile(srcFile, content, 0644); err != nil {
		t.Fatal(err)
	}

	if err := CopyFile(srcFile, dstFile); err != nil {
		t.Errorf("CopyFile() error = %v", err)
	}

	copied, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(copied) != string(content) {
		t.Errorf("CopyFile() content = %q, want %q", string(copied), string(content))
	}

	// Test copying non-existent file
	if err := CopyFile(filepath.Join(tmpDir, "nonexistent"), dstFile); err == nil {
		t.Error("CopyFile() expected error for non-existent source")
	}
}

func TestReadFileString(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "test content"

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	got := ReadFileString(testFile)
	if got != content {
		t.Errorf("ReadFileString() = %q, want %q", got, content)
	}

	// Test non-existent file
	got = ReadFileString(filepath.Join(tmpDir, "nonexistent"))
	if got != "" {
		t.Errorf("ReadFileString() for non-existent = %q, want empty", got)
	}
}

func TestGetExtension(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/path/to/file.go", ".go"},
		{"/path/to/file.test.js", ".js"},
		{"file.rs", ".rs"},
		{"/path/to/file", ""},
		{"noextension", ""},
		{"/path.dir/file", ""},
		{"", ""},
		{".hidden", ".hidden"},
	}

	for _, tt := range tests {
		got := GetExtension(tt.path)
		if got != tt.want {
			t.Errorf("GetExtension(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}
