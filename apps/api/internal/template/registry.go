package template

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Manifest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Language    string `json:"language"`
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
		return nil, fmt.Errorf("Read templates directory: %w", err)
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
	path := filepath.Join(r.root, id, "manifest.json")
	b, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, fmt.Errorf("template %q not found: %w", id, err)
	}
	var m Manifest
	if err := json.Unmarshal(b, &m); err != nil {
		return Manifest{}, fmt.Errorf("Invalid manifest: %q: %w", id, err)
	}
	if m.ID == "" {
		m.ID = id
	}
	return m, nil
}

func (r *Registry) Dir(id string) string {
	return filepath.Join(r.root, id)
}
