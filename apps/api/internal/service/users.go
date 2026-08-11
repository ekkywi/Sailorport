package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
	"github.com/ekkywi/sailorport/apps/api/internal/store"
)

var validUserRoles = map[string]struct{}{
	"admin":     {},
	"developer": {},
	"viewer":    {},
}

type UserAdminRepository interface {
	List(ctx context.Context) ([]model.User, error)
	GetByID(ctx context.Context, id string) (model.User, error)
	UpdateRole(ctx context.Context, id, role string) (model.User, error)
}

type Users struct {
	repo UserAdminRepository
}

func NewUsers(repo UserAdminRepository) *Users {
	return &Users{repo: repo}
}

func (u *Users) List(ctx context.Context) ([]model.User, error) {
	users, err := u.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return users, nil
}

func (u *Users) UpdateRole(ctx context.Context, actorID, targetID string, role string) (model.User, error) {
	actorID = strings.TrimSpace(actorID)
	targetID = strings.TrimSpace(targetID)
	role = strings.TrimSpace(role)

	if actorID == "" || targetID == "" {
		return model.User{}, fmt.Errorf("%w: user id is required", ErrInvalid)
	}
	if actorID == targetID {
		return model.User{}, fmt.Errorf("%w: cannot change your own role", ErrForbidden)
	}
	if _, ok := validUserRoles[role]; !ok {
		return model.User{}, fmt.Errorf("%w: role must be admin, developer, or viewer", ErrInvalid)
	}

	out, err := u.repo.UpdateRole(ctx, targetID, role)
	if err != nil {
		return model.User{}, mapUserAdminErr(err)
	}
	return out, nil
}

func mapUserAdminErr(err error) error {
	if errors.Is(err, store.ErrNotFound) {
		return ErrNotFound
	}
	return err
}
