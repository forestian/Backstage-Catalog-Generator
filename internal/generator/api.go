package generator

import "github.com/example/backstage-catalog-generator/internal/model"

type apiEntity struct {
	APIVersion string  `yaml:"apiVersion"`
	Kind       string  `yaml:"kind"`
	Metadata   apiMeta `yaml:"metadata"`
	Spec       apiSpec `yaml:"spec"`
}

type apiMeta struct {
	Name        string `yaml:"name"`
	Title       string `yaml:"title,omitempty"`
	Description string `yaml:"description,omitempty"`
}

type apiSpec struct {
	Type       string `yaml:"type"`
	Lifecycle  string `yaml:"lifecycle"`
	Owner      string `yaml:"owner"`
	System     string `yaml:"system,omitempty"`
	Definition string `yaml:"definition"`
}

func buildAPI(a model.APIConfig) apiEntity {
	definition := "# Definition not provided"
	if a.DefinitionPath != "" {
		definition = "# See: " + a.DefinitionPath
	}
	return apiEntity{
		APIVersion: "backstage.io/v1alpha1",
		Kind:       "API",
		Metadata: apiMeta{
			Name:  a.Name,
			Title: a.Title,
		},
		Spec: apiSpec{
			Type:       a.Type,
			Lifecycle:  a.Lifecycle,
			Owner:      a.Owner,
			System:     a.System,
			Definition: definition,
		},
	}
}
