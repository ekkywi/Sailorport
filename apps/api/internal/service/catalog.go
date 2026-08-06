package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
	"github.com/ekkywi/sailorport/apps/api/internal/store"
)

var (
	ErrInvalid  = errors.New("invalid input")
	ErrNotFound = errors.New("service not found")
	ErrConflict = errors.New("service already exists")
)

// Repository is the persistence port for catalog services.
type Repository interface {
	List(ctx context.Context) ([]model.Service, error)
	Get(ctx context.Context, id string) (model.Service, error)
	Create(ctx context.Context, req model.CreateServiceRequest) (model.Service, error)
	Update(ctx context.Context, id string, req model.UpdateServiceRequest) (model.Service, error)
	Delete(ctx context.Context, id string) error
}

type Catalog struct {
	repo Repository
}

func NewCatalog(repo Repository) *Catalog {
	return &Catalog{repo: repo}
}

func (c *Catalog) List(ctx context.Context) ([]model.Service, error) {
	services, err := c.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list services: %w", err)
	}
	return services, nil
}

func (c *Catalog) Get(ctx context.Context, id string) (model.Service, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return model.Service{}, fmt.Errorf("%w: id is required", ErrInvalid)
	}
	svc, err := c.repo.Get(ctx, id)
	if err != nil {
		return model.Service{}, mapRepoErr(err)
	}
	return svc, nil
}

func (c *Catalog) Create(ctx context.Context, req model.CreateServiceRequest) (model.Service, error) {
	req, err := normalizeCreate(req)
	if err != nil {
		return model.Service{}, err
	}
	svc, err := c.repo.Create(ctx, req)
	if err != nil {
		return model.Service{}, mapRepoErr(err)
	}
	return svc, nil
}

func (c *Catalog) Update(ctx context.Context, id string, req model.UpdateServiceRequest) (model.Service, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return model.Service{}, fmt.Errorf("%w: id is required", ErrInvalid)
	}
	req, err := normalizeUpdate(req)
	if err != nil {
		return model.Service{}, err
	}
	svc, err := c.repo.Update(ctx, id, req)
	if err != nil {
		return model.Service{}, mapRepoErr(err)
	}
	return svc, nil
}

func (c *Catalog) Delete(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("%w: id is required", ErrInvalid)
	}
	if err := c.repo.Delete(ctx, id); err != nil {
		return mapRepoErr(err)
	}
	return nil
}

func normalizeCreate(req model.CreateServiceRequest) (model.CreateServiceRequest, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	req.Owner = strings.TrimSpace(req.Owner)
	req.TemplateID = strings.TrimSpace(req.TemplateID)
	req.WorkspacePath = strings.TrimSpace(req.WorkspacePath)
	if req.Name == "" {
		return req, fmt.Errorf("%w: name is required", ErrInvalid)
	}
	return req, nil
}

func normalizeUpdate(req model.UpdateServiceRequest) (model.UpdateServiceRequest, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	req.Owner = strings.TrimSpace(req.Owner)
	if req.Name == "" {
		return req, fmt.Errorf("%w: name is required", ErrInvalid)
	}
	return req, nil
}

func mapRepoErr(err error) error {
	if errors.Is(err, store.ErrNotFound) {
		return ErrNotFound
	}
	if errors.Is(err, store.ErrConflict) {
		return ErrConflict
	}
	return err
}
