package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.0.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows version number of esp-fs-tool",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("esp-fs-tool Version: %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
