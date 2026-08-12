package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
	"github.com/ekkywi/sailorport/apps/api/internal/store"
)

type fakeUserAdminRepo struct {
	users []model.User
}

func (f *fakeUserAdminRepo) List(ctx context.Context) ([]model.User, error) {
	return f.users, nil
}

func (f *fakeUserAdminRepo) GetByID(ctx context.Context, id string) (model.User, error) {
	for _, u := range f.users {
		if u.ID == id {
			return u, nil
		}
	}
	return model.User{}, ErrNotFound
}

func (f *fakeUserAdminRepo) Create(ctx context.Context, email, name, passwordHash, role string) (model.User, error) {
	for _, u := range f.users {
		if u.Email == email {
			return model.User{}, store.ErrConflict
		}
	}
	u := model.User{
		ID:    "new",
		Email: email,
		Name:  name,
		Role:  role,
	}
	f.users = append(f.users, u)
	_ = passwordHash
	return u, nil
}

func (f *fakeUserAdminRepo) UpdateRole(ctx context.Context, id, role string) (model.User, error) {
	for i, u := range f.users {
		if u.ID == id {
			f.users[i].Role = role
			return f.users[i], nil
		}
	}
	return model.User{}, ErrNotFound
}

func TestUsers_UpdateRole_forbiddenSelfChange(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{
		users: []model.User{{ID: "u1", Email: "a@x.com", Role: "admin"}},
	})
	_, err := svc.UpdateRole(context.Background(), "u1", "u1", "viewer")
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestUsers_UpdateRole_invalidRole(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{
		users: []model.User{{ID: "u1", Role: "admin"}, {ID: "u2", Role: "developer"}},
	})
	_, err := svc.UpdateRole(context.Background(), "u1", "u2", "qa")
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}
