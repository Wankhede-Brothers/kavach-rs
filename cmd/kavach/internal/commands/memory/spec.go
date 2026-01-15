package memory

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/claude/shared/pkg/toon"
	"github.com/claude/shared/pkg/util"
	"github.com/spf13/cobra"
)

var specName string

var specCmd = &cobra.Command{
	Use:   "spec",
	Short: "Spec injector",
	Long: `[SPEC]
desc: Load and inject spec files from memory bank
path: ~/.local/shared/shared-ai/memory/specs/
formats: .toon

[FLAGS]
-n, --name: Spec name to load (without extension)

[USAGE]
kavach memory spec                  # List available specs
kavach memory spec -n api-design    # Inject api-design spec

[OUTPUT]
List:   [SPECS] with available spec names
Inject: TOON-marshaled spec content
Error:  Spec not found`,
	Run: runSpecCmd,
}

func init() {
	specCmd.Flags().StringVarP(&specName, "name", "n", "", "Spec name")
}

func runSpecCmd(cmd *cobra.Command, args []string) {
	if specName == "" {
		listSpecs()
		return
	}
	injectSpec(specName)
}

func listSpecs() {
	specsDir := util.MemoryBankPath("specs")
	entries, err := os.ReadDir(specsDir)
	if err != nil {
		fmt.Println("[SPECS]")
		fmt.Println("count: 0")
		return
	}

	fmt.Println("[SPECS]")
	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			ext := filepath.Ext(name)
			if ext == ".toon" {
				fmt.Printf("- %s\n", name[:len(name)-len(ext)])
			}
		}
	}
}

func injectSpec(name string) {
	specsDir := util.MemoryBankPath("specs")
	path := filepath.Join(specsDir, name+".toon")

	if !util.FileExists(path) {
		fmt.Println("[ERROR]")
		fmt.Printf("spec: %s.toon not found\n", name)
		os.Exit(1)
	}

	bank := toon.NewMemoryBank()
	doc, err := bank.LoadFile(path)
	if err != nil {
		fmt.Println("[ERROR]")
		fmt.Printf("load: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Print(toon.Marshal(doc))
}
