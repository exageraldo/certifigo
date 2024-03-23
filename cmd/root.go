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

	//go:embed _defaults/suacunarc.toml
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
	if err := viper.ReadConfig(bytes.NewReader(defaultConfig)); err != nil {
		fmt.Println("Can't read default config:", err)
		os.Exit(1)
	}

	if cfgFile != "" {
		_, err := os.Stat(cfgFile)
		if os.IsNotExist(err) {
			fmt.Println("Config file does not exist:", cfgFile)
			os.Exit(1)
		}

		viper.SetConfigFile(cfgFile)
		if err := viper.MergeInConfig(); err != nil {
			fmt.Println("Can't merge config:", err)
			os.Exit(1)
		}

		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Can't read config:", err)
			os.Exit(1)
		}
	}

	fmt.Println(viper.AllSettings())
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
		fmt.Println(err)
		os.Exit(1)
	}
}
