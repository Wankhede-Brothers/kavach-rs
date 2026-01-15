package path

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// PathConfig represents resolved cross-platform paths
type PathConfig struct {
	SharedMemory string `json:"shared_memory"`
	SharedConfig string `json:"shared_config"`
	SharedState  string `json:"shared_state"`
	SharedCache  string `json:"shared_cache"`
	LockDir      string `json:"lock_dir"`
}

var (
	config      *PathConfig
	configMutex sync.RWMutex
)

// ResolvePaths resolves all cross-platform paths
func ResolvePaths() *PathConfig {
	configMutex.Lock()
	defer configMutex.Unlock()

	homeDir := os.Getenv("HOME")

	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("LOCALAPPDATA")
		config = &PathConfig{
			SharedMemory: filepath.Join(appData, "SharedAI", "memory"),
			SharedConfig: filepath.Join(appData, "SharedAI", "config"),
			SharedState:  filepath.Join(appData, "SharedAI", "state"),
			SharedCache:  filepath.Join(appData, "SharedAI", "cache"),
			LockDir:      filepath.Join(appData, "SharedAI", "locks"),
		}

	case "darwin":
		config = &PathConfig{
			SharedMemory: filepath.Join(homeDir, "Library", "Application Support", "SharedAI", "memory"),
			SharedConfig: filepath.Join(homeDir, "Library", "Application Support", "SharedAI", "config"),
			SharedState:  filepath.Join(homeDir, "Library", "Application Support", "SharedAI", "state"),
			SharedCache:  filepath.Join(homeDir, "Library", "Application Support", "SharedAI", "cache"),
			LockDir:      filepath.Join(homeDir, "Library", "Application Support", "SharedAI", "locks"),
		}

	default: // Linux, BSD, etc.
		xdgData := os.Getenv("XDG_DATA_HOME")
		if xdgData == "" {
			xdgData = filepath.Join(homeDir, ".local", "share")
		}
		config = &PathConfig{
			SharedMemory: filepath.Join(xdgData, "shared-ai", "memory"),
			SharedConfig: filepath.Join(xdgData, "shared-ai", "config"),
			SharedState:  filepath.Join(xdgData, "shared-ai", "state"),
			SharedCache:  filepath.Join(xdgData, "shared-ai", "cache"),
			LockDir:      filepath.Join(xdgData, "shared-ai", "locks"),
		}
	}

	for _, dir := range []string{config.SharedMemory, config.SharedConfig, config.SharedState, config.SharedCache, config.LockDir} {
		os.MkdirAll(dir, 0755)
	}

	return config
}

// GetSharedMemoryPath convenience function
func GetSharedMemoryPath() string {
	cfg := ResolvePaths()
	return cfg.SharedMemory
}
