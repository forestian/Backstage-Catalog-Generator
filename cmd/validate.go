package cmd

import (
	"fmt"
	"os"

	"github.com/example/backstage-catalog-generator/internal/parser"
	"github.com/example/backstage-catalog-generator/internal/report"
	"github.com/example/backstage-catalog-generator/internal/validate"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a services.yaml file",
	RunE:  runValidate,
}

var validateFile string

func init() {
	validateCmd.Flags().StringVar(&validateFile, "file", "", "path to services.yaml (required)")
	validateCmd.MarkFlagRequired("file")
}

func runValidate(cmd *cobra.Command, args []string) error {
	cf, err := parser.Parse(validateFile, "unknown", "experimental", "default-system")
	if err != nil {
		return err
	}

	res := validate.Validate(cf)
	report.PrintResult(os.Stdout, res)

	if res.HasErrors() {
		fmt.Fprintf(os.Stderr, "\nValidation failed with %d error(s).\n", len(res.Errors))
		os.Exit(1)
	}

	if len(res.Warnings) == 0 {
		fmt.Println("Validation passed.")
	} else {
		fmt.Printf("\nValidation passed with %d warning(s).\n", len(res.Warnings))
	}
	return nil
}
