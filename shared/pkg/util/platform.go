// Package util provides platform-aware path utilities.
package util

import (
	"os"
	"path/filepath"
	"runtime"
)

// Platform represents the operating system.
type Platform string

const (
	PlatformLinux   Platform = "linux"
	PlatformMacOS   Platform = "macos"
	PlatformWindows Platform = "windows"
	PlatformUnknown Platform = "unknown"
)

// CLIType represents the CLI tool being used.
type CLIType string

const (
	CLIClaudeCode CLIType = "claude-code"
	CLIOpenCode   CLIType = "opencode"
)

// DetectPlatform returns the current operating system.
func DetectPlatform() Platform {
	switch runtime.GOOS {
	case "linux":
		return PlatformLinux
	case "darwin":
		return PlatformMacOS
	case "windows":
		return PlatformWindows
	default:
		return PlatformUnknown
	}
}

// PlatformPaths holds all paths for a specific platform/CLI combination.
type PlatformPaths struct {
	Platform Platform
	CLI      CLIType
	Binary   string
	Memory   string
	Config   string
	Settings string
	Logs     string
	Cache    string
}

// GetPaths returns platform-specific paths for the given CLI type.
func GetPaths(cli CLIType) *PlatformPaths {
	platform := DetectPlatform()
	home := HomeDir()

	switch platform {
	case PlatformLinux:
		return getLinuxPaths(home, cli)
	case PlatformMacOS:
		return getMacOSPaths(home, cli)
	case PlatformWindows:
		return getWindowsPaths(cli)
	default:
		return getLinuxPaths(home, cli) // Default to Linux
	}
}

// getLinuxPaths returns XDG-compliant paths for Linux.
func getLinuxPaths(home string, cli CLIType) *PlatformPaths {
	// XDG Base Directory Specification
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		dataHome = filepath.Join(home, ".local", "share")
	}

	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		configHome = filepath.Join(home, ".config")
	}

	stateHome := os.Getenv("XDG_STATE_HOME")
	if stateHome == "" {
		stateHome = filepath.Join(home, ".local", "state")
	}

	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome == "" {
		cacheHome = filepath.Join(home, ".cache")
	}

	cliName := string(cli)

	paths := &PlatformPaths{
		Platform: PlatformLinux,
		CLI:      cli,
		Memory:   filepath.Join(dataHome, cliName, "memory"),
		Config:   filepath.Join(configHome, cliName),
		Logs:     filepath.Join(stateHome, cliName, "logs"),
		Cache:    filepath.Join(cacheHome, cliName),
	}

	// Binary and settings paths differ between CLIs
	if cli == CLIClaudeCode {
		paths.Binary = filepath.Join(home, ".claude", "bin", "kavach")
		paths.Settings = filepath.Join(home, ".claude", "settings.json")
	} else {
		paths.Binary = filepath.Join(home, ".local", "bin", "kavach")
		paths.Settings = filepath.Join(configHome, cliName, "settings.json")
	}

	return paths
}

// getMacOSPaths returns standard macOS paths.
func getMacOSPaths(home string, cli CLIType) *PlatformPaths {
	cliName := string(cli)
	appSupport := filepath.Join(home, "Library", "Application Support", cliName)

	paths := &PlatformPaths{
		Platform: PlatformMacOS,
		CLI:      cli,
		Memory:   filepath.Join(appSupport, "memory"),
		Config:   appSupport,
		Logs:     filepath.Join(home, "Library", "Logs", cliName),
		Cache:    filepath.Join(home, "Library", "Caches", cliName),
	}

	if cli == CLIClaudeCode {
		paths.Binary = filepath.Join(home, ".claude", "bin", "kavach")
		paths.Settings = filepath.Join(home, ".claude", "settings.json")
	} else {
		paths.Binary = "/usr/local/bin/kavach"
		paths.Settings = filepath.Join(appSupport, "settings.json")
	}

	return paths
}

// getWindowsPaths returns standard Windows paths.
func getWindowsPaths(cli CLIType) *PlatformPaths {
	appData := os.Getenv("APPDATA")
	localAppData := os.Getenv("LOCALAPPDATA")
	userProfile := os.Getenv("USERPROFILE")

	if appData == "" {
		appData = filepath.Join(userProfile, "AppData", "Roaming")
	}
	if localAppData == "" {
		localAppData = filepath.Join(userProfile, "AppData", "Local")
	}

	cliName := string(cli)

	paths := &PlatformPaths{
		Platform: PlatformWindows,
		CLI:      cli,
		Memory:   filepath.Join(appData, cliName, "memory"),
		Config:   filepath.Join(appData, cliName),
		Logs:     filepath.Join(localAppData, cliName, "logs"),
		Cache:    filepath.Join(localAppData, cliName, "cache"),
	}

	if cli == CLIClaudeCode {
		paths.Binary = filepath.Join(userProfile, ".claude", "bin", "kavach.exe")
		paths.Settings = filepath.Join(userProfile, ".claude", "settings.json")
	} else {
		paths.Binary = filepath.Join(localAppData, "Programs", "opencode", "kavach.exe")
		paths.Settings = filepath.Join(appData, cliName, "settings.json")
	}

	return paths
}

// EnsureDirectories creates all necessary directories for the CLI.
func (p *PlatformPaths) EnsureDirectories() error {
	dirs := []string{
		p.Memory,
		p.Config,
		p.Logs,
		p.Cache,
		filepath.Dir(p.Binary),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

// EnsureMemoryBank creates memory bank directory structure.
func (p *PlatformPaths) EnsureMemoryBank() error {
	categories := []string{
		"decisions",
		"graph",
		"kanban",
		"patterns",
		"proposals",
		"research",
		"roadmaps",
		"STM",
	}

	for _, cat := range categories {
		dir := filepath.Join(p.Memory, cat)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

// GetMemoryCategory returns path to a memory category.
func (p *PlatformPaths) GetMemoryCategory(category string) string {
	return filepath.Join(p.Memory, category)
}

// GetProjectMemory returns project-specific memory path.
func (p *PlatformPaths) GetProjectMemory(project, category string) string {
	return filepath.Join(p.Memory, category, project)
}
