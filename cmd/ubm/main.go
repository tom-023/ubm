package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tom-023/ubm/internal/config"
	"github.com/tom-023/ubm/internal/storage"
)

var (
	version = "1.0.0"
	store   *storage.Storage
)

func main() {
	if err := config.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing config: %v\n", err)
		os.Exit(1)
	}

	configDir, err := config.GetConfigDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting config directory: %v\n", err)
		os.Exit(1)
	}

	store, err = storage.New(configDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing storage: %v\n", err)
		os.Exit(1)
	}

	rootCmd := &cobra.Command{
		Use:   "ubm",
		Short: "URL Bookmark Manager - Interactive command-line bookmark manager",
		Long: `ubm is a command-line URL bookmark manager with interactive directory navigation.
It allows you to organize your bookmarks in a tree-like structure and access them quickly.`,
		Version: version,
	}

	rootCmd.AddCommand(
		addCmd(),
		listCmd(),
		showCmd(),
		categoryCmd(),
		moveCmd(),
		deleteCmd(),
		editCmd(),
		// importCmd(),
		// exportCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}