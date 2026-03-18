package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "roster",
	Short: "Scaffold agent personas and orchestration for AI coding harnesses",
}

func Execute() error {
	return rootCmd.Execute()
}
