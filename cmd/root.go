package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "skaffold-buildkit-wrapper",
	Short: "buildctl wrapper for skaffold custom build",
	Long:  `buildctl wrapper for skaffold custom build`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
}
