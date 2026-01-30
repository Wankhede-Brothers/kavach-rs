// Package telemetry provides the kavach telemetry report command.
package telemetry

import (
	"fmt"

	"github.com/claude/shared/pkg/telemetry"
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Show hook execution telemetry report",
	Run:   runReport,
}

// Register adds telemetry commands to the parent.
func Register(parent *cobra.Command) {
	telCmd := &cobra.Command{
		Use:   "telemetry",
		Short: "Hook telemetry and observability",
	}
	telCmd.AddCommand(reportCmd)
	parent.AddCommand(telCmd)
}

func runReport(cmd *cobra.Command, args []string) {
	report, err := telemetry.GenerateReport()
	if err != nil {
		fmt.Printf("[TELEMETRY] No data: %v\n", err)
		return
	}

	fmt.Printf("=== Hook Telemetry Report ===\n")
	fmt.Printf("Total spans: %d\n", report.TotalSpans)
	fmt.Printf("Total duration: %dms\n", report.TotalDuration)
	fmt.Println()

	fmt.Println("--- By Hook ---")
	for name, agg := range report.ByHook {
		avg := int64(0)
		if agg.Count > 0 {
			avg = agg.TotalMs / int64(agg.Count)
		}
		fmt.Printf("  %-12s count=%d avg=%dms max=%dms tokens=%d\n",
			name, agg.Count, avg, agg.MaxMs, agg.TotalToken)
	}
	fmt.Println()

	fmt.Println("--- By Result ---")
	for result, count := range report.ByResult {
		fmt.Printf("  %-10s %d\n", result, count)
	}
	fmt.Println()

	fmt.Println("--- Slowest Hooks ---")
	for _, s := range report.Slowest {
		fmt.Printf("  %s/%s %dms result=%s\n", s.Hook, s.Tool, s.DurationMs, s.Result)
	}
}
