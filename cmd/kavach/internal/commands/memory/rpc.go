package memory

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/claude/shared/pkg/patterns"
	"github.com/claude/shared/pkg/util"
	"github.com/spf13/cobra"
)

var rpcCmd = &cobra.Command{
	Use:   "rpc",
	Short: "Memory RPC operations",
	Long: `[RPC]
desc: JSON-RPC style memory operations
purpose: Programmatic memory access from hooks

[METHODS]
get:    Retrieve memory entry by category+key
set:    Store memory entry
list:   List entries in category
health: Check RPC service health

[INPUT]
stdin: {"method":"get","params":{"category":"decisions","key":"D001"}}

[USAGE]
echo '{"method":"health"}' | kavach memory rpc
echo '{"method":"list","params":{"category":"patterns"}}' | kavach memory rpc

[OUTPUT]
{"success":true,"result":{...}}
{"success":false,"error":"..."}`,
	Run: runRPCCmd,
}

// RPCRequest represents an RPC memory request.
type RPCRequest struct {
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

// RPCResponse represents an RPC memory response.
type RPCResponse struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func runRPCCmd(cmd *cobra.Command, args []string) {
	// Read JSON-RPC request from stdin
	var req RPCRequest
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		outputRPCError("invalid request: " + err.Error())
		return
	}

	switch req.Method {
	case "get":
		handleGet(req.Params)
	case "set":
		handleSet(req.Params)
	case "list":
		handleList(req.Params)
	case "health":
		handleHealth()
	default:
		outputRPCError("unknown method: " + req.Method)
	}
}

func handleGet(params map[string]interface{}) {
	category, _ := params["category"].(string)
	key, _ := params["key"].(string)

	if category == "" || key == "" {
		outputRPCError("category and key required")
		return
	}

	// P0 SECURITY: Validate category and key to prevent path traversal
	if err := patterns.ValidateIdentifier(category); err != nil {
		outputRPCError("invalid category: " + err.Error())
		return
	}
	if err := patterns.ValidateIdentifier(key); err != nil {
		outputRPCError("invalid key: " + err.Error())
		return
	}

	memDir := util.MemoryDir()

	// Try different file patterns
	pathPatterns := []string{
		filepath.Join(memDir, category, key+".toon"),
		filepath.Join(memDir, category, key+".md"),
		filepath.Join(memDir, category, key, key+".toon"),
		filepath.Join(memDir, category, "global", key+".toon"),
	}

	for _, path := range pathPatterns {
		if data, err := os.ReadFile(path); err == nil {
			outputRPCSuccess(map[string]interface{}{
				"category": category,
				"key":      key,
				"path":     path,
				"content":  string(data),
				"status":   "found",
			})
			return
		}
	}

	outputRPCError("entry not found: " + category + "/" + key)
}

func handleSet(params map[string]interface{}) {
	category, _ := params["category"].(string)
	key, _ := params["key"].(string)
	value, _ := params["value"].(string)

	if category == "" || key == "" {
		outputRPCError("category and key required")
		return
	}

	// P0 SECURITY: Validate category and key to prevent path traversal
	if err := patterns.ValidateIdentifier(category); err != nil {
		outputRPCError("invalid category: " + err.Error())
		return
	}
	if err := patterns.ValidateIdentifier(key); err != nil {
		outputRPCError("invalid key: " + err.Error())
		return
	}

	memDir := util.MemoryDir()
	categoryDir := filepath.Join(memDir, category)

	// Ensure category directory exists
	if err := os.MkdirAll(categoryDir, 0755); err != nil {
		outputRPCError("failed to create category: " + err.Error())
		return
	}

	// Determine file extension
	ext := ".toon"
	if strings.HasSuffix(key, ".md") {
		ext = ""
	}

	filePath := filepath.Join(categoryDir, key+ext)

	// Write content
	if err := os.WriteFile(filePath, []byte(value), 0644); err != nil {
		outputRPCError("failed to write: " + err.Error())
		return
	}

	outputRPCSuccess(map[string]interface{}{
		"category": category,
		"key":      key,
		"path":     filePath,
		"status":   "stored",
	})
}

func handleList(params map[string]interface{}) {
	category, _ := params["category"].(string)

	// P0 FIX: Implement actual listing from memory bank
	memDir := util.MemoryDir()

	var entries []map[string]string

	if category == "" {
		// List all categories
		dirs, err := os.ReadDir(memDir)
		if err != nil {
			outputRPCError("failed to read memory: " + err.Error())
			return
		}

		for _, d := range dirs {
			if d.IsDir() {
				entries = append(entries, map[string]string{
					"name": d.Name(),
					"type": "category",
				})
			}
		}
	} else {
		// List entries in category
		categoryDir := filepath.Join(memDir, category)
		files, err := os.ReadDir(categoryDir)
		if err != nil {
			outputRPCError("category not found: " + category)
			return
		}

		for _, f := range files {
			entryType := "file"
			if f.IsDir() {
				entryType = "directory"
			}
			entries = append(entries, map[string]string{
				"name": f.Name(),
				"type": entryType,
			})
		}
	}

	outputRPCSuccess(map[string]interface{}{
		"category": category,
		"count":    len(entries),
		"entries":  entries,
	})
}

func handleHealth() {
	outputRPCSuccess(map[string]string{
		"status": "healthy",
	})
}

func outputRPCSuccess(result interface{}) {
	resp := RPCResponse{Success: true, Result: result}
	json.NewEncoder(os.Stdout).Encode(resp)
}

func outputRPCError(msg string) {
	resp := RPCResponse{Success: false, Error: msg}
	json.NewEncoder(os.Stdout).Encode(resp)
}
