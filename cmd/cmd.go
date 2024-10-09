package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "anime_downloader",
	Short: "Download anime from the command line",
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(setPathCmd)
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
