package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/example/backstage-catalog-generator/internal/model"
)

func sampleCatalog() *model.CatalogFile {
	return &model.CatalogFile{
		Systems: []model.SystemConfig{
			{Name: "payment", Title: "Payment System", Owner: "platform", Domain: "commerce"},
		},
		Services: []model.ServiceConfig{
			{
				Name:      "payment-api",
				Title:     "Payment API",
				Type:      "service",
				Lifecycle: "production",
				Owner:     "payment-team",
				System:    "payment",
				Repo:      "https://github.com/example/payment-api",
				Tags:      []string{"go", "api"},
				Annotations: map[string]string{
					"github.com/project-slug": "example/payment-api",
				},
				ProvidesAPIs: []string{"payment-api"},
				DependsOn:    []string{"resource:payment-db"},
			},
			{
				Name:      "worker-service",
				Title:     "Worker Service",
				Type:      "service",
				Lifecycle: "production",
				Owner:     "payment-team",
				System:    "payment",
				Repo:      "https://github.com/example/worker-service",
				Tags:      []string{"go", "worker"},
			},
		},
		APIs: []model.APIConfig{
			{Name: "payment-api", Title: "Payment API", Type: "openapi", Lifecycle: "production", Owner: "payment-team", System: "payment", DefinitionPath: "./openapi/payment.yaml"},
		},
		Resources: []model.ResourceConfig{
			{Name: "payment-db", Title: "Payment Database", Type: "database", Owner: "payment-team", System: "payment"},
			{Name: "redis-cache", Title: "Redis Cache", Type: "cache", Owner: "platform", System: "payment"},
		},
	}
}

func TestGenerate_FilesFormat(t *testing.T) {
	dir := t.TempDir()
	opts := Options{Format: "files", IncludeLocation: true, Force: false}
	files, err := Generate(sampleCatalog(), dir, opts)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	relPaths := map[string]bool{}
	for _, f := range files {
		relPaths[filepath.ToSlash(f.RelPath)] = true
	}

	required := []string{
		"payment-api/catalog-info.yaml",
		"worker-service/catalog-info.yaml",
		"systems/payment-system.yaml",
		"resources/payment-db.yaml",
		"resources/redis-cache.yaml",
		"apis/payment-api-openapi.yaml",
		"locations.yaml",
	}
	for _, r := range required {
		if !relPaths[r] {
			t.Errorf("missing expected file: %s", r)
		}
	}
}

func TestGenerate_SingleFormat(t *testing.T) {
	dir := t.TempDir()
	opts := Options{Format: "single", IncludeLocation: true, Force: false}
	files, err := Generate(sampleCatalog(), dir, opts)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	relPaths := map[string]bool{}
	for _, f := range files {
		relPaths[filepath.ToSlash(f.RelPath)] = true
	}

	if !relPaths["catalog-info.yaml"] {
		t.Error("missing catalog-info.yaml in single format")
	}
	if !relPaths["locations.yaml"] {
		t.Error("missing locations.yaml in single format")
	}
}

func TestGenerate_SingleFileContainsMultipleDocs(t *testing.T) {
	dir := t.TempDir()
	opts := Options{Format: "single", IncludeLocation: false, Force: false}
	files, err := Generate(sampleCatalog(), dir, opts)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	var content []byte
	for _, f := range files {
		if filepath.ToSlash(f.RelPath) == "catalog-info.yaml" {
			content = f.Content
		}
	}
	if content == nil {
		t.Fatal("catalog-info.yaml not found")
	}
	if !strings.Contains(string(content), "---") {
		t.Error("single file should contain YAML document separator ---")
	}
}

func TestGenerate_ComponentYAML(t *testing.T) {
	dir := t.TempDir()
	opts := Options{Format: "files", IncludeLocation: false, Force: false}
	files, err := Generate(sampleCatalog(), dir, opts)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	var content []byte
	for _, f := range files {
		if filepath.ToSlash(f.RelPath) == "payment-api/catalog-info.yaml" {
			content = f.Content
		}
	}
	if content == nil {
		t.Fatal("payment-api/catalog-info.yaml not found")
	}
	s := string(content)
	if !strings.Contains(s, "kind: Component") {
		t.Error("component YAML missing kind: Component")
	}
	if !strings.Contains(s, "apiVersion: backstage.io/v1alpha1") {
		t.Error("component YAML missing apiVersion")
	}
}

func TestGenerate_SystemYAML(t *testing.T) {
	dir := t.TempDir()
	opts := Options{Format: "files", IncludeLocation: false, Force: false}
	files, err := Generate(sampleCatalog(), dir, opts)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	var content []byte
	for _, f := range files {
		if filepath.ToSlash(f.RelPath) == "systems/payment-system.yaml" {
			content = f.Content
		}
	}
	if content == nil {
		t.Fatal("systems/payment-system.yaml not found")
	}
	if !strings.Contains(string(content), "kind: System") {
		t.Error("system YAML missing kind: System")
	}
}

func TestGenerate_APIYAML(t *testing.T) {
	dir := t.TempDir()
	opts := Options{Format: "files", IncludeLocation: false, Force: false}
	files, err := Generate(sampleCatalog(), dir, opts)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	var content []byte
	for _, f := range files {
		if filepath.ToSlash(f.RelPath) == "apis/payment-api-openapi.yaml" {
			content = f.Content
		}
	}
	if content == nil {
		t.Fatal("apis/payment-api-openapi.yaml not found")
	}
	if !strings.Contains(string(content), "kind: API") {
		t.Error("api YAML missing kind: API")
	}
}

func TestGenerate_ResourceYAML(t *testing.T) {
	dir := t.TempDir()
	opts := Options{Format: "files", IncludeLocation: false, Force: false}
	files, err := Generate(sampleCatalog(), dir, opts)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	var content []byte
	for _, f := range files {
		if filepath.ToSlash(f.RelPath) == "resources/payment-db.yaml" {
			content = f.Content
		}
	}
	if content == nil {
		t.Fatal("resources/payment-db.yaml not found")
	}
	if !strings.Contains(string(content), "kind: Resource") {
		t.Error("resource YAML missing kind: Resource")
	}
}

func TestGenerate_LocationYAML(t *testing.T) {
	dir := t.TempDir()
	opts := Options{Format: "files", IncludeLocation: true, Force: false}
	files, err := Generate(sampleCatalog(), dir, opts)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	var content []byte
	for _, f := range files {
		if filepath.ToSlash(f.RelPath) == "locations.yaml" {
			content = f.Content
		}
	}
	if content == nil {
		t.Fatal("locations.yaml not found")
	}
	s := string(content)
	if !strings.Contains(s, "kind: Location") {
		t.Error("location YAML missing kind: Location")
	}
	if !strings.Contains(s, "payment-api/catalog-info.yaml") {
		t.Error("location missing payment-api target")
	}
}

func TestGenerate_OverwriteProtection(t *testing.T) {
	dir := t.TempDir()
	opts := Options{Format: "files", IncludeLocation: false, Force: false}

	// First generate
	_, err := Generate(sampleCatalog(), dir, opts)
	if err != nil {
		t.Fatalf("first generate error: %v", err)
	}

	// Second generate without --force should fail
	_, err = Generate(sampleCatalog(), dir, opts)
	if err == nil {
		t.Error("expected error when overwriting without --force")
	}
}

func TestGenerate_ForceOverwrite(t *testing.T) {
	dir := t.TempDir()
	opts := Options{Format: "files", IncludeLocation: false, Force: false}

	_, err := Generate(sampleCatalog(), dir, opts)
	if err != nil {
		t.Fatalf("first generate error: %v", err)
	}

	opts.Force = true
	_, err = Generate(sampleCatalog(), dir, opts)
	if err != nil {
		t.Errorf("force overwrite failed: %v", err)
	}
}

func TestGenerate_PhysicalFilesCreated(t *testing.T) {
	dir := t.TempDir()
	opts := Options{Format: "files", IncludeLocation: true, Force: false}
	_, err := Generate(sampleCatalog(), dir, opts)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	expected := []string{
		filepath.Join(dir, "payment-api", "catalog-info.yaml"),
		filepath.Join(dir, "worker-service", "catalog-info.yaml"),
		filepath.Join(dir, "systems", "payment-system.yaml"),
		filepath.Join(dir, "resources", "payment-db.yaml"),
		filepath.Join(dir, "locations.yaml"),
	}
	for _, p := range expected {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			t.Errorf("expected file not created: %s", p)
		}
	}
}
