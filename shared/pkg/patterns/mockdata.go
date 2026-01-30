// Package patterns provides detection logic for mock/hardcoded data.
// mockdata.go: Blocks mock arrays, fake names, hardcoded JSON responses.
package patterns

import (
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// Frontend: const mock/dummy/fake/sample/placeholder variable names
	reMockConst = regexp.MustCompile(`(?i)\bconst\s+(mock|dummy|fake|sample|placeholder)\w*\s*[=:]`)

	// Frontend: hardcoded array with 3+ objects containing id fields
	reHardcodedArray = regexp.MustCompile(`(?s)\[\s*\{[^}]*\bid\s*:.*?\}\s*,\s*\{[^}]*\bid\s*:.*?\}\s*,\s*\{[^}]*\bid\s*:`)

	// Frontend: useState initialized with hardcoded object array
	reUseStateArray = regexp.MustCompile(`useState\(\s*\[\s*\{`)

	// Frontend: fake engagement numbers pattern
	reFakeEngagement = regexp.MustCompile(`(?i)(likes|followers|posts|views|subscribers)\s*:\s*\d{3,}`)

	// Frontend: fake names in const arrays
	reFakeNames = regexp.MustCompile(`(?i)['"](?:Pastor David|Sarah Johnson|John Smith|Jane Doe|Bob Wilson|Mary Williams)['"]`)

	// Backend Rust: vec![serde_json::json!({ in handlers
	reVecJSON = regexp.MustCompile(`vec!\s*\[\s*(?:serde_json::)?json!\s*\(`)

	// Backend: NULL as distance in SQL
	reNullDistance = regexp.MustCompile(`(?i)NULL\s+as\s+distance`)

	// Backend Rust: todo!()/unimplemented!() — checked contextually
	reTodoMacro = regexp.MustCompile(`(?:todo|unimplemented)!\s*\(`)
)

// DetectMockData checks content for mock/hardcoded data patterns.
// Returns (detected bool, reason string).
func DetectMockData(filePath, content string) (bool, string) {
	if content == "" {
		return false, ""
	}

	// Skip allowlisted files
	if isAllowlisted(filePath) {
		return false, ""
	}

	if isFrontendFile(filePath) {
		return detectFrontend(filePath, content)
	}
	if isBackendFile(filePath) {
		return detectBackend(filePath, content)
	}
	return false, ""
}

func detectFrontend(_ string, content string) (bool, string) {
	if m := reMockConst.FindString(content); m != "" {
		return true, "frontend_mock_const:" + strings.TrimSpace(m)
	}
	if reHardcodedArray.MatchString(content) {
		return true, "frontend_hardcoded_array:3+_objects_with_id_fields"
	}
	if reUseStateArray.MatchString(content) {
		return true, "frontend_useState_hardcoded_array"
	}
	if m := reFakeEngagement.FindString(content); m != "" {
		return true, "frontend_fake_engagement:" + m
	}
	if m := reFakeNames.FindString(content); m != "" {
		return true, "frontend_fake_name:" + m
	}
	return false, ""
}

func detectBackend(_ string, content string) (bool, string) {
	if reVecJSON.MatchString(content) {
		return true, "backend_hardcoded_json_vec"
	}
	if reNullDistance.MatchString(content) {
		return true, "backend_null_as_distance"
	}
	if reTodoMacro.MatchString(content) {
		// Allow in test functions
		if !strings.Contains(content, "#[test]") && !strings.Contains(content, "#[cfg(test)]") {
			return true, "backend_todo_unimplemented_in_handler"
		}
	}
	return false, ""
}

func isAllowlisted(path string) bool {
	return isTestFile(path) || isMigrationFile(path) || isDeferredPlatform(path) || isStorybook(path) || isExampleFile(path)
}

func isFrontendFile(path string) bool {
	p := strings.ToLower(path)
	return strings.HasSuffix(p, ".tsx") || strings.HasSuffix(p, ".ts") ||
		strings.HasSuffix(p, ".jsx") || strings.HasSuffix(p, ".js") ||
		strings.HasSuffix(p, ".astro")
}

func isAstroFile(path string) bool {
	return strings.HasSuffix(strings.ToLower(path), ".astro")
}

func isRustFile(path string) bool {
	return strings.HasSuffix(strings.ToLower(path), ".rs")
}

func isBackendFile(path string) bool {
	p := strings.ToLower(path)
	return strings.HasSuffix(p, ".rs") || strings.HasSuffix(p, ".go") ||
		strings.HasSuffix(p, ".py") || strings.HasSuffix(p, ".java") ||
		strings.HasSuffix(p, ".kt")
}

func isGoFile(path string) bool {
	return strings.HasSuffix(strings.ToLower(path), ".go")
}

func isPythonFile(path string) bool {
	return strings.HasSuffix(strings.ToLower(path), ".py")
}

func isJavaFile(path string) bool {
	p := strings.ToLower(path)
	return strings.HasSuffix(p, ".java") || strings.HasSuffix(p, ".kt")
}

func isDockerfile(path string) bool {
	base := strings.ToLower(filepath.Base(path))
	return base == "dockerfile" || strings.HasSuffix(base, ".dockerfile")
}

func isShellFile(path string) bool {
	p := strings.ToLower(path)
	return strings.HasSuffix(p, ".sh") || strings.HasSuffix(p, ".bash") || strings.HasSuffix(p, ".zsh")
}

func isTestFile(path string) bool {
	p := strings.ToLower(path)
	return strings.Contains(p, "_test.go") || strings.Contains(p, "test_") ||
		strings.HasSuffix(p, ".test.ts") || strings.HasSuffix(p, ".test.tsx") ||
		strings.HasSuffix(p, ".spec.ts") || strings.HasSuffix(p, ".spec.tsx") ||
		strings.Contains(p, "/tests/") || strings.Contains(p, "/test/")
}

func isMigrationFile(path string) bool {
	p := strings.ToLower(path)
	return strings.Contains(p, "/migrations/") || strings.Contains(p, "/seeds/")
}

func isDeferredPlatform(path string) bool {
	p := strings.ToLower(path)
	return strings.Contains(p, "black-diamond-fire") || strings.Contains(p, "bdf") ||
		strings.Contains(p, "xo-media") || strings.Contains(p, "xo_media")
}

func isStorybook(path string) bool {
	return strings.HasSuffix(strings.ToLower(path), ".stories.tsx") ||
		strings.HasSuffix(strings.ToLower(path), ".stories.ts")
}

func isExampleFile(path string) bool {
	p := strings.ToLower(path)
	return strings.Contains(p, "/examples/") || strings.Contains(p, "/docs/")
}
