package cmd

import (
	"fmt"
	"os"

	"github.com/exageraldo/suacuna-cli/config"
	"github.com/spf13/cobra"
)

func init() {
	var cfgFile string
	cobra.OnInitialize(func() {
		if err := config.LoadCertificateDefaults(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		if cfgFile != "" {
			if err := config.LoadCertificateFromFile(cfgFile); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
		}
	})
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}

var rootCmd = &cobra.Command{
	Use:   "suacuna-cli",
	Short: "suacuna-cli is a CLI tool to generate certificates for events.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
