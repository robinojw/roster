package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	commit    = "none"
	buildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print roster version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("roster %s (%s) built %s\n", version, commit, buildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
