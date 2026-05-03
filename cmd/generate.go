package cmd

import (
	"fmt"
	"os"

	"github.com/example/backstage-catalog-generator/internal/generator"
	"github.com/example/backstage-catalog-generator/internal/parser"
	"github.com/example/backstage-catalog-generator/internal/report"
	"github.com/example/backstage-catalog-generator/internal/validate"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Backstage catalog YAML files from services.yaml",
	RunE:  runGenerate,
}

var (
	genFile            string
	genOutput          string
	genFormat          string
	genOwner           string
	genSystem          string
	genLifecycle       string
	genIncludeLocation bool
	genForce           bool
)

func init() {
	generateCmd.Flags().StringVar(&genFile, "file", "", "path to services.yaml (required)")
	generateCmd.Flags().StringVar(&genOutput, "output", "./catalog", "output directory")
	generateCmd.Flags().StringVar(&genFormat, "format", "files", "output format: files or single")
	generateCmd.Flags().StringVar(&genOwner, "owner", "unknown", "default owner")
	generateCmd.Flags().StringVar(&genSystem, "system", "default-system", "default system")
	generateCmd.Flags().StringVar(&genLifecycle, "lifecycle", "experimental", "default lifecycle")
	generateCmd.Flags().BoolVar(&genIncludeLocation, "include-location", true, "generate locations.yaml")
	generateCmd.Flags().BoolVar(&genForce, "force", false, "overwrite existing output files")
	generateCmd.MarkFlagRequired("file")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	if genFormat != "files" && genFormat != "single" {
		return fmt.Errorf("invalid format %q: must be 'files' or 'single'", genFormat)
	}
	if genOutput == "" {
		return fmt.Errorf("output must not be empty")
	}

	cf, err := parser.Parse(genFile, genOwner, genLifecycle, genSystem)
	if err != nil {
		return err
	}

	res := validate.Validate(cf)
	report.PrintResult(os.Stdout, res)
	if res.HasErrors() {
		return fmt.Errorf("validation failed; fix errors before generating")
	}

	opts := generator.Options{
		Format:          genFormat,
		IncludeLocation: genIncludeLocation,
		Force:           genForce,
	}

	files, err := generator.Generate(cf, genOutput, opts)
	if err != nil {
		return err
	}

	fmt.Printf("Generated %d file(s) in %s\n", len(files), genOutput)
	for _, f := range files {
		fmt.Printf("  %s\n", f.RelPath)
	}
	return nil
}
