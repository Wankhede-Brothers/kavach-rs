// Package commands provides CLI commands for kavach.
// tools.go: Show Rust/Zig CLI tool status and module versions.
// DACE: Single responsibility - tool status only.
// P3 FIX: Added module versioning support.
package commands

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// Tool specification
type toolSpec struct {
	legacy  string
	rust    string
	alias   string
	purpose string
}

var rustTools = []toolSpec{
	{"cat", "bat", "cat", "syntax highlighting"},
	{"ls", "eza", "ls, ll, lt", "icons + git status"},
	{"find", "fd", "find", "10x faster search"},
	{"grep", "rg", "grep", "ripgrep (fastest)"},
	{"sed", "sd", "sed", "simpler syntax"},
	{"ps", "procs", "ps, psa", "colorful + tree"},
	{"du", "dust", "du", "visual disk usage"},
	{"top", "btm", "top, btm", "GPU graphs"},
	{"diff", "delta", "diff", "git-aware diff"},
}

var zigTools = []toolSpec{
	{"node", "bun", "bun", "4x faster JS/TS"},
	{"gnome-terminal", "ghostty", "ghostty", "GPU terminal"},
}

var toolsShowVersions bool

var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "Show Rust/Zig CLI tool status",
	Long: `[TOOLS]
desc: Display DACE-optimized CLI tool status
stack: Rust + Zig

[RUST_CLI]
cat → bat:    Syntax highlighting
ls → eza:     Icons + git status
find → fd:    10x faster search
grep → rg:    Ripgrep (fastest)
sed → sd:     Simpler syntax
ps → procs:   Colorful + tree
du → dust:    Visual disk usage
top → btm:    GPU graphs
diff → delta: Git-aware diff

[ZIG_CLI]
node → bun:      4x faster JS/TS
terminal → ghostty: GPU accelerated

[INSTALL]
Linux:   cargo install bat eza fd-find ripgrep sd procs du-dust bottom git-delta
macOS:   brew install bat eza fd ripgrep sd procs dust bottom git-delta
Windows: scoop install bat eza fd ripgrep sd procs dust bottom delta`,
	Run: runToolsCmd,
}

func init() {
	toolsCmd.Flags().BoolVar(&toolsShowVersions, "versions", false, "Show tool versions")
	rootCmd.AddCommand(toolsCmd)
}

func runToolsCmd(cmd *cobra.Command, args []string) {
	// P3 FIX: Support --versions flag
	if toolsShowVersions {
		fmt.Println("[RUST_CLI_VERSIONS]")
		fmt.Printf("%-8s %-12s %s\n", "Tool", "Version", "Status")
		fmt.Printf("%-8s %-12s %s\n", "----", "-------", "------")
		for _, t := range rustTools {
			status := checkTool(t.rust)
			version := "-"
			if status == "OK" {
				version = getToolVersion(t.rust)
			}
			fmt.Printf("%-8s %-12s %s\n", t.rust, version, status)
		}
		fmt.Println()
		fmt.Println("[ZIG_CLI_VERSIONS]")
		fmt.Printf("%-8s %-12s %s\n", "Tool", "Version", "Status")
		fmt.Printf("%-8s %-12s %s\n", "----", "-------", "------")
		for _, t := range zigTools {
			status := checkTool(t.rust)
			version := "-"
			if status == "OK" {
				version = getToolVersion(t.rust)
			}
			fmt.Printf("%-8s %-12s %s\n", t.rust, version, status)
		}
		return
	}

	fmt.Println("[RUST_CLI]")
	fmt.Printf("%-8s %-8s %-12s %-10s %s\n", "Legacy", "Rust", "Alias", "Status", "Purpose")
	fmt.Printf("%-8s %-8s %-12s %-10s %s\n", "------", "----", "-----", "------", "-------")

	for _, t := range rustTools {
		status := checkTool(t.rust)
		fmt.Printf("%-8s %-8s %-12s %-10s %s\n", t.legacy, t.rust, t.alias, status, t.purpose)
	}

	fmt.Println()
	fmt.Println("[ZIG_CLI]")
	fmt.Printf("%-8s %-8s %-12s %-10s %s\n", "Legacy", "Zig", "Alias", "Status", "Purpose")
	fmt.Printf("%-8s %-8s %-12s %-10s %s\n", "------", "---", "-----", "------", "-------")

	for _, t := range zigTools {
		status := checkTool(t.rust)
		fmt.Printf("%-8s %-8s %-12s %-10s %s\n", t.legacy, t.rust, t.alias, status, t.purpose)
	}

	fmt.Println()
	fmt.Println("[INSTALL]")
	printInstallCmd()
}

func checkTool(name string) string {
	_, err := exec.LookPath(name)
	if err != nil {
		return "missing"
	}
	return "OK"
}

// P3 FIX: Get tool version
func getToolVersion(name string) string {
	versionFlags := map[string]string{
		"bat":     "--version",
		"eza":     "--version",
		"fd":      "--version",
		"rg":      "--version",
		"sd":      "--version",
		"procs":   "--version",
		"dust":    "--version",
		"btm":     "--version",
		"delta":   "--version",
		"bun":     "--version",
		"ghostty": "--version",
	}

	flag, ok := versionFlags[name]
	if !ok {
		return "unknown"
	}

	cmd := exec.Command(name, flag)
	out, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	// Parse version from output (usually first line)
	version := strings.TrimSpace(string(out))
	lines := strings.Split(version, "\n")
	if len(lines) > 0 {
		// Extract version number
		parts := strings.Fields(lines[0])
		for _, p := range parts {
			if strings.Contains(p, ".") && (strings.HasPrefix(p, "v") || strings.HasPrefix(p, "0") || strings.HasPrefix(p, "1") || strings.HasPrefix(p, "2")) {
				return strings.TrimPrefix(p, "v")
			}
		}
		if len(parts) > 1 {
			return parts[len(parts)-1]
		}
	}
	return version
}

func printInstallCmd() {
	switch runtime.GOOS {
	case "linux":
		fmt.Println("cargo install bat eza fd-find ripgrep sd procs du-dust bottom git-delta")
		fmt.Println("# Or use kavach install script:")
		fmt.Println("curl -fsSL https://raw.githubusercontent.com/.../install.sh | bash -s -- --rust-cli")
	case "darwin":
		fmt.Println("brew install bat eza fd ripgrep sd procs dust bottom git-delta")
	case "windows":
		fmt.Println("scoop install bat eza fd ripgrep sd procs dust bottom delta")
	default:
		fmt.Println("# Install via cargo:")
		fmt.Println("cargo install bat eza fd-find ripgrep sd procs du-dust bottom git-delta")
	}
}
