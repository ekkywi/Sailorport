package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
	tpl "github.com/ekkywi/sailorport/apps/api/internal/template"
)

var namePattern = regexp.MustCompile(`^[a-z][a-z0-9-]{1,46}[a-z0-9]$`)

type ScaffoldRequest struct {
	TemplateID  string `json:"template_id"`
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	Description string `json:"description"`
}

type ScaffoldResult struct {
	Service model.Service `json:"service"`
}

type Scaffold struct {
	catalog      *Catalog
	templates    *tpl.Registry
	workspaceDir string
}

func NewScaffold(catalog *Catalog, templates *tpl.Registry, workspaceDir string) *Scaffold {
	return &Scaffold{
		catalog:      catalog,
		templates:    templates,
		workspaceDir: workspaceDir,
	}
}

func (s *Scaffold) ListTemplates() ([]tpl.Manifest, error) {
	return s.templates.List()
}

func (s *Scaffold) Run(ctx context.Context, req ScaffoldRequest) (ScaffoldResult, error) {
	req.TemplateID = strings.TrimSpace(req.TemplateID)
	req.Name = strings.TrimSpace(req.Name)
	req.Owner = strings.TrimSpace(req.Owner)
	req.Description = strings.TrimSpace(req.Description)

	if req.TemplateID == "" {
		return ScaffoldResult{}, fmt.Errorf("%w: template_id is required", ErrInvalid)
	}
	if !namePattern.MatchString(req.Name) {
		return ScaffoldResult{}, fmt.Errorf("%w: name must be lowercase kebab-case (e.g. payments-api)", ErrInvalid)
	}

	manifest, err := s.templates.Get(req.TemplateID)
	if err != nil {
		return ScaffoldResult{}, fmt.Errorf("%w: unknown template", ErrInvalid)
	}

	if err := os.MkdirAll(s.workspaceDir, 0o755); err != nil {
		return ScaffoldResult{}, fmt.Errorf(
			"create workspace dir: %w (dir not writable — see startup logs / chown data/ or use Compose named volume)",
			err,
		)
	}

	target := filepath.Join(s.workspaceDir, req.Name)
	cleanWorkspace, err := filepath.Abs(filepath.Clean(s.workspaceDir))
	if err != nil {
		return ScaffoldResult{}, fmt.Errorf("resolve workspace: %w", err)
	}
	cleanTarget, err := filepath.Abs(filepath.Clean(target))
	if err != nil {
		return ScaffoldResult{}, fmt.Errorf("resolve target: %w", err)
	}
	sep := string(filepath.Separator)
	if cleanTarget != cleanWorkspace && !strings.HasPrefix(cleanTarget+sep, cleanWorkspace+sep) {
		return ScaffoldResult{}, fmt.Errorf("%w: invalid target path", ErrInvalid)
	}
	if _, err := os.Stat(cleanTarget); err == nil {
		return ScaffoldResult{}, fmt.Errorf("%w: workspace folder already exists", ErrConflict)
	} else if !os.IsNotExist(err) {
		return ScaffoldResult{}, fmt.Errorf("stat target: %w", err)
	}

	if req.Description == "" {
		req.Description = manifest.Description
	}

	if err := tpl.Generate(tpl.GenerateInput{
		TemplateDir: s.templates.Dir(req.TemplateID),
		TargetDir:   cleanTarget,
		Name:        req.Name,
	}); err != nil {
		return ScaffoldResult{}, fmt.Errorf("generate template: %w", err)
	}

	svc, err := s.catalog.Create(ctx, model.CreateServiceRequest{
		Name:          req.Name,
		Description:   req.Description,
		Owner:         req.Owner,
		TemplateID:    manifest.ID,
		WorkspacePath: cleanTarget,
	})
	if err != nil {
		return ScaffoldResult{}, err
	}

	return ScaffoldResult{Service: svc}, nil
}
