package parser

import (
	"fmt"
	"os"

	"github.com/example/backstage-catalog-generator/internal/model"
	"github.com/example/backstage-catalog-generator/internal/normalize"
	"gopkg.in/yaml.v3"
)

// Parse reads and parses a services.yaml file, applying defaults from global config
// and CLI overrides.
func Parse(path string, cliOwner, cliLifecycle, cliSystem string) (*model.CatalogFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", path, err)
	}

	var cf model.CatalogFile
	if err := yaml.Unmarshal(data, &cf); err != nil {
		return nil, fmt.Errorf("parsing YAML in %s: %w", path, err)
	}

	applyDefaults(&cf, cliOwner, cliLifecycle, cliSystem)
	return &cf, nil
}

func applyDefaults(cf *model.CatalogFile, cliOwner, cliLifecycle, cliSystem string) {
	g := &cf.Global

	if g.Owner == "" {
		g.Owner = cliOwner
	}
	if g.Lifecycle == "" {
		g.Lifecycle = cliLifecycle
	}
	if g.System == "" {
		g.System = cliSystem
	}
	if g.Namespace == "" {
		g.Namespace = "default"
	}

	for i := range cf.Services {
		s := &cf.Services[i]
		if s.Owner == "" {
			s.Owner = g.Owner
		}
		if s.Lifecycle == "" {
			s.Lifecycle = g.Lifecycle
		}
		if s.System == "" {
			s.System = g.System
		}
		if s.Type == "" {
			s.Type = "service"
		}
		if s.Title == "" {
			s.Title = s.Name
		}
		if s.Annotations == nil {
			s.Annotations = map[string]string{}
		}
		if s.Tags == nil {
			s.Tags = []string{}
		}
		if s.ProvidesAPIs == nil {
			s.ProvidesAPIs = []string{}
		}
		if s.ConsumesAPIs == nil {
			s.ConsumesAPIs = []string{}
		}
		if s.DependsOn == nil {
			s.DependsOn = []string{}
		}
		s.Name = normalize.EntityName(s.Name)
	}

	for i := range cf.APIs {
		a := &cf.APIs[i]
		if a.Owner == "" {
			a.Owner = g.Owner
		}
		if a.Lifecycle == "" {
			a.Lifecycle = g.Lifecycle
		}
		if a.System == "" {
			a.System = g.System
		}
		if a.Type == "" {
			a.Type = "openapi"
		}
		if a.Title == "" {
			a.Title = a.Name
		}
		a.Name = normalize.EntityName(a.Name)
	}

	for i := range cf.Resources {
		r := &cf.Resources[i]
		if r.Owner == "" {
			r.Owner = g.Owner
		}
		if r.System == "" {
			r.System = g.System
		}
		if r.Type == "" {
			r.Type = "other"
		}
		if r.Title == "" {
			r.Title = r.Name
		}
		r.Name = normalize.EntityName(r.Name)
	}

	for i := range cf.Systems {
		sys := &cf.Systems[i]
		if sys.Owner == "" {
			sys.Owner = g.Owner
		}
		if sys.Title == "" {
			sys.Title = sys.Name
		}
		sys.Name = normalize.EntityName(sys.Name)
	}
}
