package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Current Git tag or the name of the snapshot
	// (https://goreleaser.com/cookbooks/using-main.version/)
	version = "dev"

	ConfigFileFromCLI string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&ConfigFileFromCLI, "config", "", "config file")
	rootCmd.AddCommand(generateCmd)
}

var rootCmd = &cobra.Command{
	Use:     "certifigo-cli",
	Short:   "certifigo-cli tool to generate certificates for events.",
	Version: version,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
