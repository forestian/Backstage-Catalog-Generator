package validate

import (
	"strings"
	"testing"

	"github.com/example/backstage-catalog-generator/internal/model"
)

func service(name, owner, lifecycle, system, repo string, tags []string) model.ServiceConfig {
	return model.ServiceConfig{
		Name:      name,
		Owner:     owner,
		Lifecycle: lifecycle,
		System:    system,
		Repo:      repo,
		Tags:      tags,
	}
}

func TestValidate_NoErrors_WithValidInput(t *testing.T) {
	cf := &model.CatalogFile{
		Services: []model.ServiceConfig{
			service("payment-api", "team", "production", "payment", "https://github.com/x/y", []string{"go"}),
		},
		Systems: []model.SystemConfig{{Name: "payment", Owner: "team"}},
	}
	res := Validate(cf)
	if res.HasErrors() {
		t.Errorf("unexpected errors: %v", res.Errors)
	}
}

func TestValidate_EmptyServiceName(t *testing.T) {
	cf := &model.CatalogFile{
		Services: []model.ServiceConfig{service("", "team", "production", "payment", "https://x.com", []string{"go"})},
	}
	res := Validate(cf)
	if !res.HasErrors() {
		t.Error("expected error for empty service name")
	}
}

func TestValidate_DuplicateServiceName(t *testing.T) {
	svc := service("payment-api", "team", "production", "payment", "https://x.com", []string{"go"})
	cf := &model.CatalogFile{Services: []model.ServiceConfig{svc, svc}}
	res := Validate(cf)
	if !res.HasErrors() {
		t.Error("expected error for duplicate service name")
	}
}

func TestValidate_DuplicateAPIName(t *testing.T) {
	api := model.APIConfig{Name: "payment-api", Owner: "team", Lifecycle: "production", Type: "openapi", DefinitionPath: "./x.yaml"}
	cf := &model.CatalogFile{APIs: []model.APIConfig{api, api}}
	res := Validate(cf)
	if !res.HasErrors() {
		t.Error("expected error for duplicate api name")
	}
}

func TestValidate_DuplicateResourceName(t *testing.T) {
	res := model.ResourceConfig{Name: "payment-db", Owner: "team", Type: "database"}
	cf := &model.CatalogFile{Resources: []model.ResourceConfig{res, res}}
	r := Validate(cf)
	if !r.HasErrors() {
		t.Error("expected error for duplicate resource name")
	}
}

func TestValidate_WarnMissingRepo(t *testing.T) {
	cf := &model.CatalogFile{
		Services: []model.ServiceConfig{service("svc", "team", "production", "sys", "", []string{"go"})},
	}
	res := Validate(cf)
	found := false
	for _, w := range res.Warnings {
		if strings.Contains(w, "missing repo") {
			found = true
		}
	}
	if !found {
		t.Error("expected warning for missing repo")
	}
}

func TestValidate_WarnHTTPRepo(t *testing.T) {
	cf := &model.CatalogFile{
		Services: []model.ServiceConfig{service("svc", "team", "production", "sys", "http://github.com/x/y", []string{"go"})},
	}
	res := Validate(cf)
	found := false
	for _, w := range res.Warnings {
		if strings.Contains(w, "plain HTTP") {
			found = true
		}
	}
	if !found {
		t.Error("expected warning for plain HTTP repo")
	}
}

func TestValidate_WarnNoServices(t *testing.T) {
	cf := &model.CatalogFile{}
	res := Validate(cf)
	found := false
	for _, w := range res.Warnings {
		if strings.Contains(w, "no services") {
			found = true
		}
	}
	if !found {
		t.Error("expected warning for no services")
	}
}

func TestValidate_WarnNoSystems(t *testing.T) {
	cf := &model.CatalogFile{
		Services: []model.ServiceConfig{service("svc", "team", "production", "sys", "https://x.com", []string{"go"})},
	}
	res := Validate(cf)
	found := false
	for _, w := range res.Warnings {
		if strings.Contains(w, "no systems") {
			found = true
		}
	}
	if !found {
		t.Error("expected warning for no systems")
	}
}

func TestValidate_WarnAPIMissingDefinitionPath(t *testing.T) {
	cf := &model.CatalogFile{
		APIs: []model.APIConfig{{Name: "my-api", Owner: "team", Lifecycle: "production", Type: "openapi"}},
	}
	res := Validate(cf)
	found := false
	for _, w := range res.Warnings {
		if strings.Contains(w, "definition_path") {
			found = true
		}
	}
	if !found {
		t.Error("expected warning for missing definition_path")
	}
}
