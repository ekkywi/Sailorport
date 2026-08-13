package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/auth"
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
	Create(ctx context.Context, email, name, passwordHash, role string) (model.User, error)
	UpdateRole(ctx context.Context, id, role string) (model.User, error)
	UpdateDisabled(ctx context.Context, id string, disabled bool) (model.User, error)
	UpdatePasswordHash(ctx context.Context, id, passwordHash string) error
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

func (u *Users) Create(ctx context.Context, req model.CreateUserRequest) (model.User, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Name = strings.TrimSpace(req.Name)
	req.Password = strings.TrimSpace(req.Password)
	req.Role = strings.TrimSpace(req.Role)

	if req.Email == "" || req.Password == "" {
		return model.User{}, fmt.Errorf("%w: email and password are required", ErrInvalid)
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return model.User{}, fmt.Errorf("%w: invalid email", ErrInvalid)
	}
	if len(req.Password) < 8 {
		return model.User{}, fmt.Errorf("%w: password must be at least 8 characters", ErrInvalid)
	}
	if req.Name == "" {
		req.Name = strings.Split(req.Email, "@")[0]
	}
	if req.Role == "" {
		req.Role = "developer"
	}
	if _, ok := validUserRoles[req.Role]; !ok {
		return model.User{}, fmt.Errorf("%w: role must be admin, developer, or viewer", ErrInvalid)
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return model.User{}, fmt.Errorf("hash password: %w", err)
	}

	out, err := u.repo.Create(ctx, req.Email, req.Name, hash, req.Role)
	if err != nil {
		return model.User{}, mapUserAdminErr(err)
	}
	return out, nil
}

func (u *Users) UpdateRole(ctx context.Context, actorID, targetID string, role string) (model.User, error) {
	actorID = strings.TrimSpace(actorID)
	targetID = strings.TrimSpace(targetID)
	role = strings.TrimSpace(role)

	if actorID == "" || targetID == "" {
		return model.User{}, fmt.Errorf("%w: user id is required", ErrInvalid)
	}
	if !isUserID(targetID) {
		return model.User{}, fmt.Errorf("%w: invalid user id", ErrInvalid)
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

func (u *Users) SetDisabled(ctx context.Context, actorID, targetID string, disabled bool) (model.User, error) {
	actorID = strings.TrimSpace(actorID)
	targetID = strings.TrimSpace(targetID)

	if actorID == "" || targetID == "" {
		return model.User{}, fmt.Errorf("%w: user id is required", ErrInvalid)
	}
	if !isUserID(targetID) {
		return model.User{}, fmt.Errorf("%w: invalid user id", ErrInvalid)
	}
	if actorID == targetID {
		return model.User{}, fmt.Errorf("%w: cannot disable or enable yourself", ErrForbidden)
	}

	out, err := u.repo.UpdateDisabled(ctx, targetID, disabled)
	if err != nil {
		return model.User{}, mapUserAdminErr(err)
	}
	return out, nil
}

func (u *Users) ResetPassword(ctx context.Context, actorID, targetID, password string) error {
	actorID = strings.TrimSpace(actorID)
	targetID = strings.TrimSpace(targetID)
	password = strings.TrimSpace(password)

	if actorID == "" || targetID == "" {
		return fmt.Errorf("%w: user id is required", ErrInvalid)
	}
	if !isUserID(targetID) {
		return fmt.Errorf("%w: invalid user id", ErrInvalid)
	}
	if actorID == targetID {
		return fmt.Errorf("%w: cannot reset your own password here", ErrForbidden)
	}
	if len(password) < 8 {
		return fmt.Errorf("%w: password must be at least 8 characters", ErrInvalid)
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	if err := u.repo.UpdatePasswordHash(ctx, targetID, hash); err != nil {
		return mapUserAdminErr(err)
	}
	return nil
}

// isUserID checks for a UUID-shaped id (8-4-4-4-12 hex) before hitting Postgres.
func isUserID(id string) bool {
	if len(id) != 36 {
		return false
	}
	for i, c := range id {
		switch i {
		case 8, 13, 18, 23:
			if c != '-' {
				return false
			}
		default:
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
				return false
			}
		}
	}
	return true
}

func mapUserAdminErr(err error) error {
	if errors.Is(err, store.ErrNotFound) {
		return ErrNotFound
	}
	if errors.Is(err, store.ErrConflict) {
		return ErrConflict
	}
	return err
}
