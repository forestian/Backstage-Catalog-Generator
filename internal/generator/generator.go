package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/example/backstage-catalog-generator/internal/model"
	"gopkg.in/yaml.v3"
)

type Options struct {
	Format          string // "files" or "single"
	IncludeLocation bool
	Force           bool
}

type GeneratedFile struct {
	RelPath string
	Content []byte
}

// Generate produces catalog YAML files from a parsed CatalogFile.
func Generate(cf *model.CatalogFile, outputDir string, opts Options) ([]GeneratedFile, error) {
	var files []GeneratedFile
	var locationTargets []string

	if opts.Format == "single" {
		content, err := buildSingleFile(cf)
		if err != nil {
			return nil, err
		}
		files = append(files, GeneratedFile{RelPath: "catalog-info.yaml", Content: content})
		locationTargets = append(locationTargets, "./catalog-info.yaml")
	} else {
		// per-service files
		for _, svc := range cf.Services {
			entity := buildComponent(svc)
			content, err := marshalYAML(entity)
			if err != nil {
				return nil, fmt.Errorf("marshaling component %s: %w", svc.Name, err)
			}
			rel := filepath.Join(svc.Name, "catalog-info.yaml")
			files = append(files, GeneratedFile{RelPath: rel, Content: content})
			locationTargets = append(locationTargets, "./"+filepath.ToSlash(rel))
		}

		// system files (deduplicated by name)
		seenSystems := map[string]bool{}
		for _, sys := range cf.Systems {
			if seenSystems[sys.Name] {
				continue
			}
			seenSystems[sys.Name] = true
			entity := buildSystem(sys)
			content, err := marshalYAML(entity)
			if err != nil {
				return nil, fmt.Errorf("marshaling system %s: %w", sys.Name, err)
			}
			rel := filepath.Join("systems", sys.Name+"-system.yaml")
			files = append(files, GeneratedFile{RelPath: rel, Content: content})
			locationTargets = append(locationTargets, "./"+filepath.ToSlash(rel))
		}

		// resource files (deduplicated by name)
		seenResources := map[string]bool{}
		for _, res := range cf.Resources {
			if seenResources[res.Name] {
				continue
			}
			seenResources[res.Name] = true
			entity := buildResource(res)
			content, err := marshalYAML(entity)
			if err != nil {
				return nil, fmt.Errorf("marshaling resource %s: %w", res.Name, err)
			}
			rel := filepath.Join("resources", res.Name+".yaml")
			files = append(files, GeneratedFile{RelPath: rel, Content: content})
			locationTargets = append(locationTargets, "./"+filepath.ToSlash(rel))
		}

		// api files (deduplicated by name)
		seenAPIs := map[string]bool{}
		for _, api := range cf.APIs {
			if seenAPIs[api.Name] {
				continue
			}
			seenAPIs[api.Name] = true
			entity := buildAPI(api)
			content, err := marshalYAML(entity)
			if err != nil {
				return nil, fmt.Errorf("marshaling api %s: %w", api.Name, err)
			}
			rel := filepath.Join("apis", api.Name+"-"+api.Type+".yaml")
			files = append(files, GeneratedFile{RelPath: rel, Content: content})
			locationTargets = append(locationTargets, "./"+filepath.ToSlash(rel))
		}
	}

	if opts.IncludeLocation && len(locationTargets) > 0 {
		loc := buildLocation(locationTargets)
		content, err := marshalYAML(loc)
		if err != nil {
			return nil, fmt.Errorf("marshaling location: %w", err)
		}
		files = append(files, GeneratedFile{RelPath: "locations.yaml", Content: content})
	}

	if err := writeFiles(files, outputDir, opts.Force); err != nil {
		return nil, err
	}
	return files, nil
}

func buildSingleFile(cf *model.CatalogFile) ([]byte, error) {
	var buf bytes.Buffer
	first := true

	appendEntity := func(v interface{}) error {
		if !first {
			buf.WriteString("---\n")
		}
		first = false
		b, err := marshalYAML(v)
		if err != nil {
			return err
		}
		buf.Write(b)
		return nil
	}

	for _, svc := range cf.Services {
		if err := appendEntity(buildComponent(svc)); err != nil {
			return nil, err
		}
	}
	for _, sys := range cf.Systems {
		if err := appendEntity(buildSystem(sys)); err != nil {
			return nil, err
		}
	}
	for _, res := range cf.Resources {
		if err := appendEntity(buildResource(res)); err != nil {
			return nil, err
		}
	}
	for _, api := range cf.APIs {
		if err := appendEntity(buildAPI(api)); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func marshalYAML(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	enc.Close()
	return buf.Bytes(), nil
}

func writeFiles(files []GeneratedFile, outputDir string, force bool) error {
	// check for conflicts first
	if !force {
		for _, f := range files {
			dest := filepath.Join(outputDir, f.RelPath)
			if _, err := os.Stat(dest); err == nil {
				return fmt.Errorf("output file already exists: %s (use --force to overwrite)", dest)
			}
		}
	}

	for _, f := range files {
		dest := filepath.Join(outputDir, f.RelPath)
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return fmt.Errorf("creating directory for %s: %w", dest, err)
		}
		if err := os.WriteFile(dest, f.Content, 0644); err != nil {
			return fmt.Errorf("writing %s: %w", dest, err)
		}
	}
	return nil
}
