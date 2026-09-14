package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
	"github.com/ekkywi/sailorport/apps/api/internal/store"
)

type fakeSetupUsers struct {
	n     int
	err   error
	users []model.User
}

func (f *fakeSetupUsers) Count(ctx context.Context) (int, error) {
	if f.err != nil {
		return 0, f.err
	}
	if len(f.users) > 0 {
		return len(f.users), nil
	}
	return f.n, nil
}

func (f *fakeSetupUsers) Create(ctx context.Context, email, name, passwordHash, role string) (model.User, error) {
	for _, u := range f.users {
		if u.Email == email {
			return model.User{}, store.ErrConflict
		}
	}
	u := model.User{
		ID:    "admin-1",
		Email: email,
		Name:  name,
		Role:  role,
	}
	f.users = append(f.users, u)
	_ = passwordHash
	return u, nil
}

func TestSetupStatus_NeedsSetupWhenEmpty(t *testing.T) {
	s := NewSetup(&fakeSetupUsers{n: 0})
	st, err := s.Status(context.Background())
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if !st.NeedsSetup {
		t.Fatalf("expected needs_setup true when count=0")
	}
}

func TestSetupStatus_ReadyWhenUsersExist(t *testing.T) {
	s := NewSetup(&fakeSetupUsers{n: 1})
	st, err := s.Status(context.Background())
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if st.NeedsSetup {
		t.Fatalf("expected needs_setup false when count>0")
	}
}

func TestSetupStatus_CountError(t *testing.T) {
	s := NewSetup(&fakeSetupUsers{err: errors.New("db down")})
	_, err := s.Status(context.Background())
	if err == nil {
		t.Fatalf("expected error when Cound fails")
	}
}

func TestSetupCreateAdmin_OK(t *testing.T) {
	s := NewSetup(&fakeSetupUsers{n: 0})
	u, err := s.CreateAdmin(context.Background(), model.SetupAdminRequest{
		Email:    "admin@example.com",
		Name:     "Admin User",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("CreateAdmin: %v", err)
	}
	if u.Email != "admin@example.com" {
		t.Fatalf("email lowercased: got %q", u.Email)
	}
	if u.Role != "admin" {
		t.Fatalf("role: got %q want admin", u.Role)
	}
	if u.Name != "Admin User" {
		t.Fatalf("name: got %q", u.Name)
	}
}

func TestSetupCreateAdmin_ClosedWhenUsersExist(t *testing.T) {
	s := NewSetup(&fakeSetupUsers{n: 1})
	_, err := s.CreateAdmin(context.Background(), model.SetupAdminRequest{
		Email:    "admin@example.com",
		Password: "password123",
	})
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestSetupCreateAdmin_InvalidEmail(t *testing.T) {
	s := NewSetup(&fakeSetupUsers{n: 0})
	_, err := s.CreateAdmin(context.Background(), model.SetupAdminRequest{
		Email:    "not-an-email",
		Password: "password123",
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}

func TestSetupCreateAdmin_ShortPassword(t *testing.T) {
	s := NewSetup(&fakeSetupUsers{n: 0})
	_, err := s.CreateAdmin(context.Background(), model.SetupAdminRequest{
		Email:    "a@example.com",
		Password: "short",
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}

func TestSetupCreateAdmin_DefaultNameFromEmail(t *testing.T) {
	s := NewSetup(&fakeSetupUsers{n: 0})
	u, err := s.CreateAdmin(context.Background(), model.SetupAdminRequest{
		Email:    "admin@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("CreateAdmin: %v", err)
	}
	if u.Name != "admin" {
		t.Fatalf("default name: got %q want admin", u.Name)
	}
}
