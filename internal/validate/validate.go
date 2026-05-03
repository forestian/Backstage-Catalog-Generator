package validate

import (
	"fmt"
	"strings"

	"github.com/example/backstage-catalog-generator/internal/model"
)

type Result struct {
	Errors   []string
	Warnings []string
}

func (r *Result) HasErrors() bool { return len(r.Errors) > 0 }

func (r *Result) addError(msg string)   { r.Errors = append(r.Errors, msg) }
func (r *Result) addWarning(msg string) { r.Warnings = append(r.Warnings, msg) }

// Validate checks a parsed CatalogFile for errors and warnings.
func Validate(cf *model.CatalogFile) *Result {
	res := &Result{}

	if len(cf.Services) == 0 {
		res.addWarning("no services defined")
	}
	if len(cf.Systems) == 0 {
		res.addWarning("no systems defined")
	}

	serviceNames := map[string]bool{}
	for _, s := range cf.Services {
		if s.Name == "" {
			res.addError("service has empty name")
			continue
		}
		if serviceNames[s.Name] {
			res.addError(fmt.Sprintf("duplicate service name: %q", s.Name))
		}
		serviceNames[s.Name] = true

		if s.Owner == "" {
			res.addWarning(fmt.Sprintf("service %q missing owner", s.Name))
		}
		if s.System == "" {
			res.addWarning(fmt.Sprintf("service %q missing system", s.Name))
		}
		if s.Repo == "" {
			res.addWarning(fmt.Sprintf("service %q missing repo", s.Name))
		} else if strings.HasPrefix(s.Repo, "http://") {
			res.addWarning(fmt.Sprintf("service %q repo uses plain HTTP instead of HTTPS", s.Name))
		}
		if len(s.Tags) == 0 {
			res.addWarning(fmt.Sprintf("service %q has no tags", s.Name))
		}
		if s.Lifecycle == "experimental" && isProductionLike(s.System) {
			res.addWarning(fmt.Sprintf("service %q has lifecycle=experimental in production-like system %q", s.Name, s.System))
		}
	}

	apiNames := map[string]bool{}
	for _, a := range cf.APIs {
		if a.Name == "" {
			res.addError("api has empty name")
			continue
		}
		if apiNames[a.Name] {
			res.addError(fmt.Sprintf("duplicate api name: %q", a.Name))
		}
		apiNames[a.Name] = true
		if a.DefinitionPath == "" {
			res.addWarning(fmt.Sprintf("api %q missing definition_path", a.Name))
		}
	}

	resourceNames := map[string]bool{}
	for _, r := range cf.Resources {
		if r.Name == "" {
			res.addError("resource has empty name")
			continue
		}
		if resourceNames[r.Name] {
			res.addError(fmt.Sprintf("duplicate resource name: %q", r.Name))
		}
		resourceNames[r.Name] = true
		if r.Type == "" {
			res.addWarning(fmt.Sprintf("resource %q missing type", r.Name))
		}
	}

	return res
}

func isProductionLike(system string) bool {
	s := strings.ToLower(system)
	return strings.Contains(s, "prod") || strings.Contains(s, "payment") || strings.Contains(s, "billing")
}
