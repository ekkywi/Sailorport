package service

import (
	"context"
	"fmt"
	"net/mail"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/auth"
	"github.com/ekkywi/sailorport/apps/api/internal/model"
)

type SetupUsers interface {
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, email, name, passwordHash, role string) (model.User, error)
}

type Setup struct {
	users SetupUsers
}

func NewSetup(users SetupUsers) *Setup {
	return &Setup{users: users}
}

func (s *Setup) Status(ctx context.Context) (model.SetupStatus, error) {
	count, err := s.users.Count(ctx)
	if err != nil {
		return model.SetupStatus{}, fmt.Errorf("count users: %w", err)
	}
	return model.SetupStatus{
		NeedsSetup: count == 0,
	}, nil
}

func (s *Setup) CreateAdmin(ctx context.Context, req model.SetupAdminRequest) (model.User, error) {
	count, err := s.users.Count(ctx)
	if err != nil {
		return model.User{}, fmt.Errorf("count users: %w", err)
	}
	if count > 0 {
		return model.User{}, fmt.Errorf(
			"%w: setup is closed — admin already exists",
			ErrForbidden,
		)
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Name = strings.TrimSpace(req.Name)
	req.Password = strings.TrimSpace(req.Password)

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

	const role = "admin"

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return model.User{}, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.users.Create(ctx, req.Email, req.Name, hash, role)
	if err != nil {
		return model.User{}, mapUserErr(err)
	}
	return user, nil
}
