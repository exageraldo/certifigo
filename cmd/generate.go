package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate certificates for events.",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}
