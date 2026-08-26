package catalogapp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Manifest struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Image         string   `json:"image"`
	ContainerPort int      `json:"container_port"`
	Tags          []string `json:"tags,omitempty"`
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
	return m, nil
}
