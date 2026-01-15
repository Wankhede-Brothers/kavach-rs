package util

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// FileExists checks if a file exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// DirExists checks if a directory exists.
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// EnsureDir creates a directory if it doesn't exist.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// EnsureParentDir creates the parent directory if it doesn't exist.
func EnsureParentDir(filePath string) error {
	return EnsureDir(filepath.Dir(filePath))
}

// ReadJSON reads and unmarshals JSON from a file.
func ReadJSON(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// WriteJSON marshals and writes JSON to a file.
func WriteJSON(path string, v interface{}) error {
	if err := EnsureParentDir(path); err != nil {
		return err
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// AppendFile appends data to a file, creating it if needed.
func AppendFile(path string, data []byte) error {
	if err := EnsureParentDir(path); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

// CopyFile copies a file from src to dst.
func CopyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := EnsureParentDir(dst); err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

// ReadFileString reads file and returns content as string.
func ReadFileString(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

// GetExtension returns file extension including the dot.
func GetExtension(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			return path[i:]
		}
		if path[i] == '/' {
			return ""
		}
	}
	return ""
}
