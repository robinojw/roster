package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/robinojw/roster/internal/analyser"
	"github.com/robinojw/roster/internal/writer"
)

var (
	bootstrapPath   string
	bootstrapDryRun bool
)

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Analyse repo and scaffold agent personas with orchestration",
	RunE:  runBootstrap,
}

func init() {
	bootstrapCmd.Flags().StringVar(&bootstrapPath, "path", ".", "repo root to analyse")
	bootstrapCmd.Flags().BoolVar(&bootstrapDryRun, "dry-run", false, "print what would be written")
	rootCmd.AddCommand(bootstrapCmd)
}

func runBootstrap(cmd *cobra.Command, args []string) error {
	signals, err := analyser.Analyse(bootstrapPath)
	if err != nil {
		return fmt.Errorf("analyse repo: %w", err)
	}

	if bootstrapDryRun {
		return printDryRun(signals)
	}

	result, err := writer.WriteAll(bootstrapPath, signals)
	if err != nil {
		return fmt.Errorf("write files: %w", err)
	}

	printSummary(result)
	return nil
}

func printDryRun(signals *analyser.RepoSignals) error {
	data, err := json.MarshalIndent(signals, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal signals: %w", err)
	}
	fmt.Println("Repo signals (dry run):")
	fmt.Println(string(data))
	fmt.Println("\nFiles that would be written:")
	fmt.Println("  .roster/signals.json")
	fmt.Println("  .roster/personas/ (11 persona files)")
	fmt.Println("  CLAUDE.md (managed section)")
	fmt.Println("  AGENTS.md (managed section)")
	return nil
}

func printSummary(result *writer.Result) {
	fmt.Printf("roster: wrote %d files\n", len(result.FilesWritten))
	for _, file := range result.FilesWritten {
		fmt.Printf("  %s\n", file)
	}
}
