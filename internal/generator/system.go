package generator

import "github.com/example/backstage-catalog-generator/internal/model"

type systemEntity struct {
	APIVersion string     `yaml:"apiVersion"`
	Kind       string     `yaml:"kind"`
	Metadata   systemMeta `yaml:"metadata"`
	Spec       systemSpec `yaml:"spec"`
}

type systemMeta struct {
	Name        string `yaml:"name"`
	Title       string `yaml:"title,omitempty"`
	Description string `yaml:"description,omitempty"`
}

type systemSpec struct {
	Owner  string `yaml:"owner"`
	Domain string `yaml:"domain,omitempty"`
}

func buildSystem(s model.SystemConfig) systemEntity {
	return systemEntity{
		APIVersion: "backstage.io/v1alpha1",
		Kind:       "System",
		Metadata: systemMeta{
			Name:        s.Name,
			Title:       s.Title,
			Description: s.Description,
		},
		Spec: systemSpec{
			Owner:  s.Owner,
			Domain: s.Domain,
		},
	}
}
