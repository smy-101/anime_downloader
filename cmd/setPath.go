package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	dirPermissions = 0755
)

var setPathCmd = &cobra.Command{
	Use:   "setPath",
	Short: "Set the output path for downloaded files",
	Args:  cobra.ExactArgs(1),
	Run:   runSetPath,
}

func runSetPath(cmd *cobra.Command, args []string) {
	path := args[0]

	// Get absolute path and clean it
	absPath, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		handleError("Error getting absolute path", err)
	}

	// Check if directory exists, create if necessary
	if err := ensureDirectoryExists(absPath); err != nil {
		handleError("Error ensuring directory exists", err)
	}

	// Check write permissions
	if err := checkWritePermissions(absPath); err != nil {
		handleError("Error checking write permissions", err)
	}

	// Save path to config
	if err := savePathToConfig(absPath); err != nil {
		handleError("Error saving path to config", err)
	}

	fmt.Printf("Successfully set default path to: %s\n", absPath)
}

func ensureDirectoryExists(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Printf("Directory does not exist: %s\n", path)
		if promptYesNo("Do you want to create this directory?") {
			return os.MkdirAll(path, dirPermissions)
		}
		fmt.Println("Directory creation cancelled. Exiting.")
		os.Exit(0)
	}
	return err
}

func checkWritePermissions(path string) error {
	testFile := filepath.Join(path, ".write_test")
	f, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("no write permission in the directory: %v", err)
	}
	f.Close()
	return os.Remove(testFile)
}

func savePathToConfig(path string) error {
	viper.Set("download_dir", path)
	return viper.WriteConfig()
}

func promptYesNo(question string) bool {
	prompt := promptui.Prompt{
		Label: question + " [y/n]",
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Println("Error prompting user:", err)
		return false
	}

	return strings.ToLower(result) == "y" || strings.ToLower(result) == "yes"
}

func handleError(message string, err error) {
	log.Fatalf("%s: %v", message, err)
}
