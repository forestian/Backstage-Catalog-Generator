package generator

type locationEntity struct {
	APIVersion string       `yaml:"apiVersion"`
	Kind       string       `yaml:"kind"`
	Metadata   locationMeta `yaml:"metadata"`
	Spec       locationSpec `yaml:"spec"`
}

type locationMeta struct {
	Name string `yaml:"name"`
}

type locationSpec struct {
	Type    string   `yaml:"type"`
	Targets []string `yaml:"targets"`
}

func buildLocation(targets []string) locationEntity {
	return locationEntity{
		APIVersion: "backstage.io/v1alpha1",
		Kind:       "Location",
		Metadata: locationMeta{
			Name: "generated-catalog-location",
		},
		Spec: locationSpec{
			Type:    "file",
			Targets: targets,
		},
	}
}
