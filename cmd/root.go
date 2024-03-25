package cmd

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	//go:embed _default.toml
	defaultConfig []byte
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	// rootCmd.PersistentFlags().StringVarP(&projectBase, "projectbase", "b", "", "base project directory eg. github.com/spf13/")
	// rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "Author name for copyright attribution")
	// rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "Name of license for the project (can provide `licensetext` in config)")
	// rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
	// viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	// viper.BindPFlag("projectbase", rootCmd.PersistentFlags().Lookup("projectbase"))
	// viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	// viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	// viper.SetDefault("license", "apache")
}

func initConfig() {
	viper.SetConfigType("toml")

	if err := viper.MergeConfig(bytes.NewReader(defaultConfig)); err != nil {
		fmt.Fprintf(os.Stderr, "Can't read default config: %s\n", err)
		os.Exit(1)
	}

	if cfgFile != "" {
		_, err := os.Stat(cfgFile)
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Config file does not exist: %s\n", err)
			os.Exit(1)
		}

		viper.SetConfigFile(cfgFile)
		if err := viper.MergeInConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Can't merge config: %s\n", err)
			os.Exit(1)
		}

		if err := viper.ReadInConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Can't read passed config: %s\n", err)
			os.Exit(1)
		}
	}
}

var rootCmd = &cobra.Command{
	Use:   "suacuna-cli",
	Short: "suacuna-cli is a CLI tool to generate certificates for events.",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
