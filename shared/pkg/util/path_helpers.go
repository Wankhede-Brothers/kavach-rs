// Package util provides common utility functions.
// path_helpers.go: Path manipulation helpers (minimal imports).
// DACE: Single responsibility - path helpers only.
package util

import "path/filepath"

// join wraps filepath.Join.
func join(elem ...string) string {
	return filepath.Join(elem...)
}

// base wraps filepath.Base.
func base(path string) string {
	return filepath.Base(path)
}

// parent wraps filepath.Dir.
func parent(path string) string {
	return filepath.Dir(path)
}

// split wraps filepath.Split.
func split(path string) (string, string) {
	return filepath.Split(path)
}

// clean wraps filepath.Clean.
func clean(path string) string {
	return filepath.Clean(path)
}
