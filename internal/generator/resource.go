package generator

import "github.com/example/backstage-catalog-generator/internal/model"

type resourceEntity struct {
	APIVersion string       `yaml:"apiVersion"`
	Kind       string       `yaml:"kind"`
	Metadata   resourceMeta `yaml:"metadata"`
	Spec       resourceSpec `yaml:"spec"`
}

type resourceMeta struct {
	Name        string `yaml:"name"`
	Title       string `yaml:"title,omitempty"`
	Description string `yaml:"description,omitempty"`
}

type resourceSpec struct {
	Type   string `yaml:"type"`
	Owner  string `yaml:"owner"`
	System string `yaml:"system,omitempty"`
}

func buildResource(r model.ResourceConfig) resourceEntity {
	return resourceEntity{
		APIVersion: "backstage.io/v1alpha1",
		Kind:       "Resource",
		Metadata: resourceMeta{
			Name:        r.Name,
			Title:       r.Title,
			Description: r.Description,
		},
		Spec: resourceSpec{
			Type:   r.Type,
			Owner:  r.Owner,
			System: r.System,
		},
	}
}
