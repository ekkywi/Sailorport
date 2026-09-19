package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
	"github.com/ekkywi/sailorport/apps/api/internal/store"
)

type fakeAuthUsers struct {
	n     int
	users []model.User
}

func (f *fakeAuthUsers) Count(ctx context.Context) (int, error) {
	if len(f.users) > 0 {
		return len(f.users), nil
	}
	return f.n, nil
}

func (f *fakeAuthUsers) Create(ctx context.Context, email, name, passwordHash, role string) (model.User, error) {
	u := model.User{
		ID:    "u1",
		Email: email,
		Name:  name,
		Role:  role,
	}
	f.users = append(f.users, u)
	_ = passwordHash
	return u, nil
}

func (f *fakeAuthUsers) GetByEmail(ctx context.Context, email string) (model.UserRecord, error) {
	return model.UserRecord{}, store.ErrNotFound
}

func (f *fakeAuthUsers) GetByID(ctx context.Context, id string) (model.User, error) {
	return model.User{}, store.ErrNotFound
}

type fakeRegGate struct {
	open bool
	err  error
}

func (f *fakeRegGate) RegistrationOpen(ctx context.Context) (bool, error) {
	if f.err != nil {
		return false, f.err
	}
	return f.open, nil
}

func TestAuthRegister_NeedsSetup(t *testing.T) {
	a := NewAuth(&fakeAuthUsers{n: 0}, &fakeRegGate{open: true}, "secret")
	_, err := a.Register(context.Background(), model.RegisterRequest{
		Email: "a@example.com", Password: "password123",
	})
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestAuthRegister_Closed(t *testing.T) {
	a := NewAuth(&fakeAuthUsers{n: 1}, &fakeRegGate{open: false}, "secret")
	_, err := a.Register(context.Background(), model.RegisterRequest{
		Email: "a@example.com", Password: "password123",
	})
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestAuthRegister_OpenCreatesDeveloper(t *testing.T) {
	users := &fakeAuthUsers{n: 1}
	a := NewAuth(users, &fakeRegGate{open: true}, "secret")
	u, err := a.Register(context.Background(), model.RegisterRequest{
		Email:    "Dev@Example.com",
		Name:     "Dev",
		Password: "password123",
		Role:     "admin", // harus diabaikan
	})
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if u.Role != "developer" {
		t.Fatalf("role: got %q want developer", u.Role)
	}
	if u.Email != "dev@example.com" {
		t.Fatalf("email: got %q", u.Email)
	}
}
