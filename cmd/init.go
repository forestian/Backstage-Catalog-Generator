package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/example/backstage-catalog-generator/internal/generator"
	"github.com/example/backstage-catalog-generator/internal/parser"
	"github.com/example/backstage-catalog-generator/internal/validate"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create an example project directory with sample services.yaml and generated catalog files",
	RunE:  runInit,
}

var initOutput string

func init() {
	initCmd.Flags().StringVar(&initOutput, "output", "./cataloggen-demo", "output directory for the demo project")
}

func runInit(cmd *cobra.Command, args []string) error {
	if _, err := os.Stat(initOutput); err == nil {
		return fmt.Errorf("output directory already exists: %s", initOutput)
	}

	if err := os.MkdirAll(initOutput, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// Write demo services.yaml
	servicesPath := filepath.Join(initOutput, "services.yaml")
	if err := os.WriteFile(servicesPath, []byte(demoServicesYAML), 0644); err != nil {
		return fmt.Errorf("writing services.yaml: %w", err)
	}

	// Parse and generate catalog
	cf, err := parser.Parse(servicesPath, "platform", "production", "payment")
	if err != nil {
		return fmt.Errorf("parsing demo services.yaml: %w", err)
	}

	res := validate.Validate(cf)
	for _, w := range res.Warnings {
		fmt.Fprintf(os.Stdout, "WARNING: %s\n", w)
	}

	catalogDir := filepath.Join(initOutput, "catalog")
	opts := generator.Options{
		Format:          "files",
		IncludeLocation: true,
		Force:           false,
	}
	files, err := generator.Generate(cf, catalogDir, opts)
	if err != nil {
		return fmt.Errorf("generating catalog: %w", err)
	}

	// Write README
	readmePath := filepath.Join(initOutput, "README.md")
	readmeContent := buildDemoREADME(initOutput, files)
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("writing README.md: %w", err)
	}

	fmt.Printf("Initialized demo project in %s\n", initOutput)
	fmt.Printf("  %s\n", "services.yaml")
	fmt.Printf("  %s\n", "README.md")
	for _, f := range files {
		fmt.Printf("  catalog/%s\n", f.RelPath)
	}
	return nil
}

func buildDemoREADME(dir string, files []generator.GeneratedFile) string {
	var buf bytes.Buffer
	buf.WriteString("# Backstage Catalog Generator — Demo\n\n")
	buf.WriteString("This directory was created by `cataloggen init`.\n\n")
	buf.WriteString("## Generated files\n\n")
	buf.WriteString("```\n")
	buf.WriteString(dir + "/\n")
	buf.WriteString("  README.md\n")
	buf.WriteString("  services.yaml\n")
	buf.WriteString("  catalog/\n")
	for _, f := range files {
		buf.WriteString("    " + f.RelPath + "\n")
	}
	buf.WriteString("```\n\n")
	buf.WriteString("## How services.yaml works\n\n")
	buf.WriteString("The `services.yaml` file describes your services, systems, APIs, and resources.\n")
	buf.WriteString("The `global` section provides defaults for owner, lifecycle, system, and namespace.\n\n")
	buf.WriteString("## How to validate\n\n")
	buf.WriteString("```sh\n")
	buf.WriteString("cataloggen validate --file ./services.yaml\n")
	buf.WriteString("```\n\n")
	buf.WriteString("## How to generate\n\n")
	buf.WriteString("```sh\n")
	buf.WriteString("cataloggen generate --file ./services.yaml --output ./catalog\n")
	buf.WriteString("cataloggen generate --file ./services.yaml --output ./catalog --force\n")
	buf.WriteString("cataloggen generate --file ./services.yaml --output ./catalog-single --format single\n")
	buf.WriteString("```\n\n")
	buf.WriteString("## How to register in Backstage\n\n")
	buf.WriteString("1. Copy the generated `catalog/` directory into your repository.\n")
	buf.WriteString("2. Register `catalog/locations.yaml` in your Backstage instance via **Register an existing component**.\n")
	buf.WriteString("3. Backstage will discover all catalog entities referenced by the Location targets.\n\n")
	buf.WriteString("## How Location targets work\n\n")
	buf.WriteString("The `locations.yaml` file is a Backstage `Location` entity. Its `spec.targets` list\n")
	buf.WriteString("points to each generated catalog-info.yaml. When registered, Backstage reads all targets.\n\n")
	buf.WriteString("## Limitations\n\n")
	buf.WriteString("- API definition files are referenced by path but not read by this tool.\n")
	buf.WriteString("- No Backstage API calls are made. Registration is manual.\n")
	buf.WriteString("- Review all generated files before registering them in Backstage.\n")
	return buf.String()
}

const demoServicesYAML = `global:
  owner: platform
  lifecycle: production
  system: payment
  namespace: default

systems:
  - name: payment
    title: Payment System
    description: Services related to payment processing
    owner: platform
    domain: commerce

services:
  - name: payment-api
    title: Payment API
    description: Handles payment requests from external clients
    type: service
    lifecycle: production
    owner: payment-team
    system: payment
    repo: https://github.com/example/payment-api
    docs: https://docs.example.com/payment-api
    tags:
      - go
      - api
      - kubernetes
    annotations:
      github.com/project-slug: example/payment-api
      backstage.io/techdocs-ref: dir:.
    provides_apis:
      - payment-api
    consumes_apis:
      - user-api
    depends_on:
      - resource:payment-db
      - resource:redis-cache

  - name: worker-service
    title: Worker Service
    description: Processes asynchronous payment jobs
    type: service
    lifecycle: production
    owner: payment-team
    system: payment
    repo: https://github.com/example/worker-service
    tags:
      - go
      - worker
      - queue
    depends_on:
      - resource:redis-cache

apis:
  - name: payment-api
    title: Payment API
    type: openapi
    lifecycle: production
    owner: payment-team
    system: payment
    definition_path: ./openapi/payment.yaml

resources:
  - name: payment-db
    title: Payment Database
    type: database
    owner: payment-team
    system: payment

  - name: redis-cache
    title: Redis Cache
    type: cache
    owner: platform
    system: payment
`
