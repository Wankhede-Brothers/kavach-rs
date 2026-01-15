// Package util provides shared utility functions.
// string.go: String manipulation utilities.
// DACE: Reusable string functions.
package util

// Itoa converts int to string without fmt package.
func Itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	if neg {
		return "-" + result
	}
	return result
}

// Min returns the smaller of two ints.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the larger of two ints.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// CountLines returns number of lines in content.
func CountLines(content string) int {
	count := 1
	for _, ch := range content {
		if ch == '\n' {
			count++
		}
	}
	return count
}
