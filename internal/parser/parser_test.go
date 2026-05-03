package parser

import (
	"os"
	"path/filepath"
	"testing"
)

const sampleYAML = `
global:
  owner: platform
  lifecycle: production
  system: payment
  namespace: default

systems:
  - name: payment
    title: Payment System

services:
  - name: payment-api
    title: Payment API
    owner: payment-team
    lifecycle: production
    system: payment
    repo: https://github.com/example/payment-api
    tags:
      - go

  - name: worker-service
    title: Worker Service

apis:
  - name: payment-api
    type: openapi

resources:
  - name: payment-db
    type: database
`

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "services.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestParse_Basic(t *testing.T) {
	path := writeTemp(t, sampleYAML)
	cf, err := Parse(path, "unknown", "experimental", "default-system")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(cf.Services) != 2 {
		t.Errorf("expected 2 services, got %d", len(cf.Services))
	}
}

func TestParse_GlobalDefaultsApplied(t *testing.T) {
	path := writeTemp(t, sampleYAML)
	cf, err := Parse(path, "cli-owner", "experimental", "cli-system")
	if err != nil {
		t.Fatal(err)
	}
	// global.owner is set in YAML, so CLI should not override it
	if cf.Global.Owner != "platform" {
		t.Errorf("global owner = %q, want %q", cf.Global.Owner, "platform")
	}
}

func TestParse_ServiceInheritsGlobalOwner(t *testing.T) {
	yaml := `
global:
  owner: platform
services:
  - name: my-svc
`
	path := writeTemp(t, yaml)
	cf, err := Parse(path, "cli-owner", "experimental", "default-system")
	if err != nil {
		t.Fatal(err)
	}
	if cf.Services[0].Owner != "platform" {
		t.Errorf("service owner = %q, want platform", cf.Services[0].Owner)
	}
}

func TestParse_ServiceInheritsGlobalLifecycle(t *testing.T) {
	yaml := `
global:
  lifecycle: production
services:
  - name: my-svc
    owner: team
`
	path := writeTemp(t, yaml)
	cf, err := Parse(path, "unknown", "experimental", "default-system")
	if err != nil {
		t.Fatal(err)
	}
	if cf.Services[0].Lifecycle != "production" {
		t.Errorf("service lifecycle = %q, want production", cf.Services[0].Lifecycle)
	}
}

func TestParse_ServiceInheritsGlobalSystem(t *testing.T) {
	yaml := `
global:
  system: my-system
services:
  - name: my-svc
    owner: team
`
	path := writeTemp(t, yaml)
	cf, err := Parse(path, "unknown", "experimental", "default-system")
	if err != nil {
		t.Fatal(err)
	}
	if cf.Services[0].System != "my-system" {
		t.Errorf("service system = %q, want my-system", cf.Services[0].System)
	}
}

func TestParse_ServiceTypeDefaultsToService(t *testing.T) {
	yaml := `
services:
  - name: my-svc
    owner: team
`
	path := writeTemp(t, yaml)
	cf, err := Parse(path, "owner", "experimental", "sys")
	if err != nil {
		t.Fatal(err)
	}
	if cf.Services[0].Type != "service" {
		t.Errorf("service type = %q, want service", cf.Services[0].Type)
	}
}

func TestParse_NameNormalized(t *testing.T) {
	yaml := `
services:
  - name: My Service
    owner: team
`
	path := writeTemp(t, yaml)
	cf, err := Parse(path, "owner", "experimental", "sys")
	if err != nil {
		t.Fatal(err)
	}
	if cf.Services[0].Name != "my-service" {
		t.Errorf("normalized name = %q, want my-service", cf.Services[0].Name)
	}
}

func TestParse_InvalidYAML(t *testing.T) {
	// unclosed flow sequence is a genuine YAML syntax error
	path := writeTemp(t, "services: [unclosed")
	_, err := Parse(path, "owner", "experimental", "sys")
	if err == nil {
		t.Error("expected parse error for invalid YAML")
	}
}

func TestParse_CLIOwnerUsedWhenNoGlobal(t *testing.T) {
	yaml := `
services:
  - name: my-svc
`
	path := writeTemp(t, yaml)
	cf, err := Parse(path, "cli-owner", "experimental", "sys")
	if err != nil {
		t.Fatal(err)
	}
	if cf.Services[0].Owner != "cli-owner" {
		t.Errorf("owner = %q, want cli-owner", cf.Services[0].Owner)
	}
}
