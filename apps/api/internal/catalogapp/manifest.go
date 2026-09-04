package catalogapp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type EnvField struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required"`
	Secret      bool   `json:"secret"`
	Default     string `json:"default,omitempty"`
}

type Version struct {
	Tag     string `json:"tag"`
	Image   string `json:"image"`
	Default bool   `json:"default,omitempty"`
}

type Manifest struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Image         string     `json:"image"`
	Versions      []Version  `json:"versions,omitempty"`
	ContainerPort int        `json:"container_port"`
	Tags          []string   `json:"tags,omitempty"`
	Env           []EnvField `json:"env,omitempty"`
	Command       []string   `json:"command,omitempty"`
}

func validateEnvFields(appID string, fields []EnvField) error {
	if len(fields) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(fields))
	for i, f := range fields {
		name := strings.TrimSpace(f.Name)
		if name == "" {
			return fmt.Errorf("catalog app %q: env[%d]: name is required", appID, i)
		}
		if _, dup := seen[name]; dup {
			return fmt.Errorf("catalog app %q: duplicate env name %q", appID, name)
		}
		seen[name] = struct{}{}

		for j, r := range name {
			if r >= 'A' && r <= 'Z' {
				continue
			}
			if r >= '0' && r <= '9' && j > 0 {
				continue
			}
			if r == '_' && j > 0 {
				continue
			}
			return fmt.Errorf("catalog app %q: invalid env name %q", appID, name)
		}
	}
	return nil
}

func validateVersions(appID string, topImage string, versions []Version) error {
	if len(versions) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(versions))
	defaultCount := 0
	var defaultImage string

	for i, v := range versions {
		tag := strings.TrimSpace(v.Tag)
		image := strings.TrimSpace(v.Image)
		if tag == "" {
			return fmt.Errorf("catalog app %q: versions[%d]: tag is required", appID, i)
		}
		if image == "" {
			return fmt.Errorf("catalog app %q: versions[%d]: image is required", appID, i)
		}
		if _, dup := seen[tag]; dup {
			return fmt.Errorf("catalog app %q: duplicate version tag %q", appID, tag)
		}
		seen[tag] = struct{}{}

		if v.Default {
			defaultCount++
			defaultImage = image
		}
	}

	if defaultCount != 1 {
		return fmt.Errorf("catalog app %q: versions require exactly one default version: true", appID)
	}

	topImage = strings.TrimSpace(topImage)
	if topImage != "" && topImage != defaultImage {
		return fmt.Errorf(
			"catalog app %q: top-level image %q must match default version image %q",
			appID, topImage, defaultImage,
		)
	}

	return nil
}

func validateCommand(appID string, env []EnvField, command []string) error {
	if len(command) == 0 {
		return nil
	}

	allowed := make(map[string]struct{}, len(env))
	for _, f := range env {
		name := strings.TrimSpace(f.Name)
		if name != "" {
			allowed[name] = struct{}{}
		}
	}

	for i, raw := range command {
		arg := strings.TrimSpace(raw)
		if arg == "" {
			return fmt.Errorf("catalog app %q: command[%d]: empty argument", appID, i)
		}

		rest := arg
		for {
			start := strings.Index(rest, "${")
			if start < 0 {
				break
			}
			end := strings.Index(rest[start:], "}")
			if end < 0 {
				return fmt.Errorf("catalog app %q: command[%d]: unclosed placeholder in %q", appID, i, arg)
			}
			end = start + end
			name := rest[start+2 : end]
			if name == "" {
				return fmt.Errorf("catalog app %q: command[%d]: empty placeholder in %q", appID, i, arg)
			}
			for j, r := range name {
				ok := (r >= 'A' && r <= 'Z') ||
					(r == '_' && j > 0) ||
					(r >= '0' && r <= '9' && j > 0)
				if !ok {
					return fmt.Errorf("catalog app %q: command[%d]: invalid placeholder %q", appID, i, name)
				}
			}
			if _, ok := allowed[name]; !ok {
				return fmt.Errorf(
					"catalog app %q: command[%d]: placeholder ${%s} is not defined in env",
					appID, i, name,
				)
			}
			rest = rest[end+1:]
		}
	}
	return nil
}

type Registry struct {
	root string
}

func NewRegistry(root string) *Registry {
	return &Registry{root: root}
}

func (r *Registry) List() ([]Manifest, error) {
	entries, err := os.ReadDir(r.root)
	if err != nil {
		if os.IsNotExist(err) {
			return []Manifest{}, nil
		}
		return nil, fmt.Errorf("read catalog-apps dir: %w", err)
	}
	out := make([]Manifest, 0)
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		m, err := r.Get(e.Name())
		if err != nil {
			continue
		}
		out = append(out, m)
	}
	return out, nil
}

func (r *Registry) Get(id string) (Manifest, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Manifest{}, fmt.Errorf("catalog app id is empty")
	}
	path := filepath.Join(r.root, id, "manifest.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, fmt.Errorf("catalog app %q not found: %w", id, err)
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return Manifest{}, fmt.Errorf("parse catalog app %q: %w", id, err)
	}
	if strings.TrimSpace(m.ID) == "" {
		m.ID = id
	}
	if strings.TrimSpace(m.Image) == "" {
		return Manifest{}, fmt.Errorf("catalog app %q: image is required", id)
	}

	for i := range m.Env {
		m.Env[i].Name = strings.TrimSpace(m.Env[i].Name)
		m.Env[i].Description = strings.TrimSpace(m.Env[i].Description)
		m.Env[i].Default = strings.TrimSpace(m.Env[i].Default)
	}

	if err := validateEnvFields(id, m.Env); err != nil {
		return Manifest{}, err
	}

	for i := range m.Versions {
		m.Versions[i].Tag = strings.TrimSpace(m.Versions[i].Tag)
		m.Versions[i].Image = strings.TrimSpace(m.Versions[i].Image)
	}

	if err := validateVersions(id, m.Image, m.Versions); err != nil {
		return Manifest{}, err
	}

	for i := range m.Command {
		m.Command[i] = strings.TrimSpace(m.Command[i])
	}
	if err := validateCommand(id, m.Env, m.Command); err != nil {
		return Manifest{}, err
	}

	return m, nil
}
