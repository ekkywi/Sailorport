package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
	"github.com/ekkywi/sailorport/apps/api/internal/store"
)

const (
	testAdminID = "00000000-0000-0000-0000-000000000001"
	testDevID   = "00000000-0000-0000-0000-000000000002"
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
	return model.User{}, store.ErrNotFound
}

func (f *fakeUserAdminRepo) UpdateDisabled(ctx context.Context, id string, disabled bool) (model.User, error) {
	for i, u := range f.users {
		if u.ID == id {
			f.users[i].Disabled = disabled
			return f.users[i], nil
		}
	}
	return model.User{}, store.ErrNotFound
}

func (f *fakeUserAdminRepo) UpdatePasswordHash(ctx context.Context, id, passwordHash string) error {
	for _, u := range f.users {
		if u.ID == id {
			_ = passwordHash
			return nil
		}
	}
	return store.ErrNotFound
}

func TestUsers_UpdateRole_forbiddenSelfChange(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{
		users: []model.User{{ID: testAdminID, Email: "a@x.com", Role: "admin"}},
	})
	_, err := svc.UpdateRole(context.Background(), testAdminID, testAdminID, "viewer")
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestUsers_UpdateRole_invalidRole(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{
		users: []model.User{
			{ID: testAdminID, Role: "admin"},
			{ID: testDevID, Role: "developer"},
		},
	})
	_, err := svc.UpdateRole(context.Background(), testAdminID, testDevID, "qa")
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}

func TestUsers_SetDisabled_forbiddenSelf(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{
		users: []model.User{{ID: testAdminID, Email: "a@x.com", Role: "admin"}},
	})
	_, err := svc.SetDisabled(context.Background(), testAdminID, testAdminID, true)
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestUsers_SetDisabled_ok(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{
		users: []model.User{
			{ID: testAdminID, Role: "admin"},
			{ID: testDevID, Role: "developer", Disabled: false},
		},
	})
	out, err := svc.SetDisabled(context.Background(), testAdminID, testDevID, true)
	if err != nil {
		t.Fatal(err)
	}
	if !out.Disabled {
		t.Fatalf("expected disabled=true")
	}
}

func TestUsers_SetDisabled_invalidID(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{
		users: []model.User{{ID: testAdminID, Role: "admin"}},
	})
	_, err := svc.SetDisabled(context.Background(), testAdminID, "TARGET_USER_ID", true)
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}

func TestUsers_ResetPassword_forbiddenSelf(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{
		users: []model.User{{ID: testAdminID, Role: "admin"}},
	})
	err := svc.ResetPassword(context.Background(), testAdminID, testAdminID, "newpassword")
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestUsers_ResetPassword_ok(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{
		users: []model.User{
			{ID: testAdminID, Role: "admin"},
			{ID: testDevID, Role: "developer"},
		},
	})
	err := svc.ResetPassword(context.Background(), testAdminID, testDevID, "newpassword")
	if err != nil {
		t.Fatal(err)
	}
}

func TestUsers_ResetPassword_tooShort(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{
		users: []model.User{
			{ID: testAdminID, Role: "admin"},
			{ID: testDevID, Role: "developer"},
		},
	})
	err := svc.ResetPassword(context.Background(), testAdminID, testDevID, "short")
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}
