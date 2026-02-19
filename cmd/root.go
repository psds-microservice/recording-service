package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "recording-service",
	Short: "Recording service: ingest stream from streaming-service, write to storage, return URL",
	Long:  `gRPC-only service. Command: api (default).`,
	RunE:  runAPI,
}

func init() {
	rootCmd.AddCommand(apiCmd)
}

// Execute runs the root command (for main to log.Fatal on error).
func Execute() error {
	return rootCmd.Execute()
}
