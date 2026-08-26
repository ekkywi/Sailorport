package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/ekkywi/sailorport/apps/api/internal/auth"
	"github.com/ekkywi/sailorport/apps/api/internal/model"
	"github.com/ekkywi/sailorport/apps/api/internal/store"
)

type UserRepository interface {
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, email, name, passwordHash, role string) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.UserRecord, error)
	GetByID(ctx context.Context, id string) (model.User, error)
}

type Auth struct {
	users     UserRepository
	jwtSecret string
	tokenTTL  time.Duration
}

func NewAuth(users UserRepository, jwtSecret string) *Auth {
	return &Auth{
		users:     users,
		jwtSecret: jwtSecret,
		tokenTTL:  24 * time.Hour,
	}
}

// Register hanya melayani bootstrap: akun pertama di instalasi baru, dan akun itu
// jadi admin. Sesudah ada user, penambahan akun lewat POST /api/v1/users (admin)
// supaya endpoint publik ini tidak bisa dipakai orang luar untuk dapat hak deploy.
func (a *Auth) Register(ctx context.Context, req model.RegisterRequest) (model.User, error) {
	count, err := a.users.Count(ctx)
	if err != nil {
		return model.User{}, fmt.Errorf("count users: %w", err)
	}
	if count > 0 {
		return model.User{}, fmt.Errorf(
			"%w: registration is closed — ask an admin to create your account",
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

	// Akun pertama memiliki platform, jadi role dari request diabaikan.
	const role = "admin"

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return model.User{}, fmt.Errorf("hash password: %w", err)
	}

	user, err := a.users.Create(ctx, req.Email, req.Name, hash, role)
	if err != nil {
		return model.User{}, mapUserErr(err)
	}
	return user, nil
}

func (a *Auth) Login(ctx context.Context, req model.LoginRequest) (model.LoginResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	password := strings.TrimSpace(req.Password)
	if email == "" || password == "" {
		return model.LoginResponse{}, fmt.Errorf("%w: email and password are required", ErrInvalid)
	}

	rec, err := a.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return model.LoginResponse{}, fmt.Errorf("%w: invalid email or password", ErrUnauthorized)
		}
		return model.LoginResponse{}, mapUserErr(err)
	}

	if !auth.CheckPassword(rec.PasswordHash, password) {
		return model.LoginResponse{}, fmt.Errorf("%w: invalid email or password", ErrUnauthorized)
	}

	if rec.Disabled {
		return model.LoginResponse{}, fmt.Errorf("%w: account is disabled", ErrUnauthorized)
	}

	token, err := auth.IssueToken(rec.ID, rec.Email, rec.Role, a.jwtSecret, a.tokenTTL)
	if err != nil {
		return model.LoginResponse{}, fmt.Errorf("issue token: %w", err)
	}

	return model.LoginResponse{
		Token: token,
		User:  rec.User,
	}, nil
}

func (a *Auth) Me(ctx context.Context, userID string) (model.User, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return model.User{}, fmt.Errorf("%w: unauthorized", ErrUnauthorized)
	}
	user, err := a.users.GetByID(ctx, userID)
	if err != nil {
		return model.User{}, mapUserErr(err)
	}
	if user.Disabled {
		return model.User{}, fmt.Errorf("%w: account is disabled", ErrUnauthorized)
	}
	return user, nil
}

func mapUserErr(err error) error {
	if errors.Is(err, store.ErrNotFound) {
		return ErrNotFound
	}
	if errors.Is(err, store.ErrConflict) {
		return ErrConflict
	}
	return err
}
