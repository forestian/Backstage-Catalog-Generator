package generator

import (
	"strings"

	"github.com/example/backstage-catalog-generator/internal/model"
)

type componentEntity struct {
	APIVersion string        `yaml:"apiVersion"`
	Kind       string        `yaml:"kind"`
	Metadata   componentMeta `yaml:"metadata"`
	Spec       componentSpec `yaml:"spec"`
}

type componentMeta struct {
	Name        string            `yaml:"name"`
	Title       string            `yaml:"title,omitempty"`
	Description string            `yaml:"description,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
	Tags        []string          `yaml:"tags,omitempty"`
	Links       []metaLink        `yaml:"links,omitempty"`
}

type metaLink struct {
	URL   string `yaml:"url"`
	Title string `yaml:"title"`
}

type componentSpec struct {
	Type         string   `yaml:"type"`
	Lifecycle    string   `yaml:"lifecycle"`
	Owner        string   `yaml:"owner"`
	System       string   `yaml:"system,omitempty"`
	ProvidesAPIs []string `yaml:"providesApis,omitempty"`
	ConsumesAPIs []string `yaml:"consumesApis,omitempty"`
	DependsOn    []string `yaml:"dependsOn,omitempty"`
}

func buildComponent(s model.ServiceConfig) componentEntity {
	annotations := map[string]string{}
	for k, v := range s.Annotations {
		annotations[k] = v
	}
	if s.Repo != "" {
		annotations["backstage.io/source-location"] = "url:" + s.Repo
		if _, ok := annotations["backstage.io/techdocs-ref"]; !ok {
			annotations["backstage.io/techdocs-ref"] = "dir:."
		}
	}
	if len(annotations) == 0 {
		annotations = nil
	}

	var tags []string
	if len(s.Tags) > 0 {
		tags = s.Tags
	}

	var links []metaLink
	if s.Docs != "" {
		links = append(links, metaLink{URL: s.Docs, Title: "Documentation"})
	}

	var providesAPIs []string
	for _, a := range s.ProvidesAPIs {
		providesAPIs = append(providesAPIs, a)
	}
	var consumesAPIs []string
	for _, a := range s.ConsumesAPIs {
		consumesAPIs = append(consumesAPIs, a)
	}
	var dependsOn []string
	for _, d := range s.DependsOn {
		dep := d
		if !strings.Contains(dep, ":") {
			dep = "resource:" + dep
		}
		dependsOn = append(dependsOn, dep)
	}

	return componentEntity{
		APIVersion: "backstage.io/v1alpha1",
		Kind:       "Component",
		Metadata: componentMeta{
			Name:        s.Name,
			Title:       s.Title,
			Description: s.Description,
			Annotations: annotations,
			Tags:        tags,
			Links:       links,
		},
		Spec: componentSpec{
			Type:         s.Type,
			Lifecycle:    s.Lifecycle,
			Owner:        s.Owner,
			System:       s.System,
			ProvidesAPIs: providesAPIs,
			ConsumesAPIs: consumesAPIs,
			DependsOn:    dependsOn,
		},
	}
}
