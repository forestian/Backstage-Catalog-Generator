package model

type CatalogFile struct {
	Global    GlobalConfig     `yaml:"global"`
	Systems   []SystemConfig   `yaml:"systems"`
	Services  []ServiceConfig  `yaml:"services"`
	APIs      []APIConfig      `yaml:"apis"`
	Resources []ResourceConfig `yaml:"resources"`
}

type GlobalConfig struct {
	Owner     string `yaml:"owner"`
	Lifecycle string `yaml:"lifecycle"`
	System    string `yaml:"system"`
	Namespace string `yaml:"namespace"`
}

type ServiceConfig struct {
	Name         string            `yaml:"name"`
	Title        string            `yaml:"title"`
	Description  string            `yaml:"description"`
	Type         string            `yaml:"type"`
	Lifecycle    string            `yaml:"lifecycle"`
	Owner        string            `yaml:"owner"`
	System       string            `yaml:"system"`
	Repo         string            `yaml:"repo"`
	Docs         string            `yaml:"docs"`
	Tags         []string          `yaml:"tags"`
	Annotations  map[string]string `yaml:"annotations"`
	ProvidesAPIs []string          `yaml:"provides_apis"`
	ConsumesAPIs []string          `yaml:"consumes_apis"`
	DependsOn    []string          `yaml:"depends_on"`
}

type SystemConfig struct {
	Name        string `yaml:"name"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Owner       string `yaml:"owner"`
	Domain      string `yaml:"domain"`
}

type APIConfig struct {
	Name           string `yaml:"name"`
	Title          string `yaml:"title"`
	Type           string `yaml:"type"`
	Lifecycle      string `yaml:"lifecycle"`
	Owner          string `yaml:"owner"`
	System         string `yaml:"system"`
	DefinitionPath string `yaml:"definition_path"`
}

type ResourceConfig struct {
	Name        string `yaml:"name"`
	Title       string `yaml:"title"`
	Type        string `yaml:"type"`
	Owner       string `yaml:"owner"`
	System      string `yaml:"system"`
	Description string `yaml:"description"`
}
