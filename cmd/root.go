package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cataloggen",
	Short: "Generate Backstage catalog-info.yaml files from a services.yaml definition",
	Long: `cataloggen reads a services.yaml file and generates Backstage catalog entities
for Components, Systems, APIs, and Resources.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(generateCmd)
}
