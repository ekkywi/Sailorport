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
	return model.User{}, store.ErrNotFound
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

func (f *fakeUserAdminRepo) SoftDelete(ctx context.Context, id string) error {
	for i, u := range f.users {
		if u.ID == id {
			f.users = append(f.users[:i], f.users[i+1:]...)
			return nil
		}
	}
	return store.ErrNotFound
}

func TestUsers_List_ok(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{
		users: []model.User{
			{ID: testAdminID, Role: "admin"},
			{ID: testDevID, Role: "developer"},
		},
	})
	out, err := svc.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 users, got %d", len(out))
	}
}

func TestUsers_Create_ok(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{})
	out, err := svc.Create(context.Background(), model.CreateUserRequest{
		Email:    "  Dev@Example.COM ",
		Password: "password1",
		Role:     "",
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.Email != "dev@example.com" {
		t.Fatalf("email: got %q", out.Email)
	}
	if out.Role != "developer" {
		t.Fatalf("role: got %q", out.Role)
	}
	if out.Name != "dev" {
		t.Fatalf("name: got %q", out.Name)
	}
}

func TestUsers_Create_withNameAndRole(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{})
	out, err := svc.Create(context.Background(), model.CreateUserRequest{
		Email:    "alice@example.com",
		Name:     "Alice",
		Password: "password1",
		Role:     "viewer",
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.Name != "Alice" || out.Role != "viewer" {
		t.Fatalf("got name=%q role=%q", out.Name, out.Role)
	}
}

func TestUsers_Create_invalidEmail(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{})
	_, err := svc.Create(context.Background(), model.CreateUserRequest{
		Email:    "not-an-email",
		Password: "password1",
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}

func TestUsers_Create_passwordTooShort(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{})
	_, err := svc.Create(context.Background(), model.CreateUserRequest{
		Email:    "dev@example.com",
		Password: "short",
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}

func TestUsers_Create_invalidRole(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{})
	_, err := svc.Create(context.Background(), model.CreateUserRequest{
		Email:    "dev@example.com",
		Password: "password1",
		Role:     "qa",
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}

func TestUsers_Create_conflict(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{
		users: []model.User{{ID: testDevID, Email: "dev@example.com"}},
	})
	_, err := svc.Create(context.Background(), model.CreateUserRequest{
		Email:    "dev@example.com",
		Password: "password1",
		Role:     "viewer",
	})
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
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

func TestUsers_UpdateRole_ok(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{
		users: []model.User{
			{ID: testAdminID, Role: "admin"},
			{ID: testDevID, Role: "developer"},
		},
	})
	out, err := svc.UpdateRole(context.Background(), testAdminID, testDevID, "viewer")
	if err != nil {
		t.Fatal(err)
	}
	if out.Role != "viewer" {
		t.Fatalf("expected viewer, got %q", out.Role)
	}
}

func TestUsers_UpdateRole_invalidID(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{
		users: []model.User{{ID: testAdminID, Role: "admin"}},
	})
	_, err := svc.UpdateRole(context.Background(), testAdminID, "TARGET_USER_ID", "viewer")
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}

func TestUsers_UpdateRole_notFound(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{
		users: []model.User{{ID: testAdminID, Role: "admin"}},
	})
	_, err := svc.UpdateRole(context.Background(), testAdminID, testDevID, "viewer")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
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

func TestUsers_SetDisabled_enable(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{
		users: []model.User{
			{ID: testAdminID, Role: "admin"},
			{ID: testDevID, Role: "developer", Disabled: true},
		},
	})
	out, err := svc.SetDisabled(context.Background(), testAdminID, testDevID, false)
	if err != nil {
		t.Fatal(err)
	}
	if out.Disabled {
		t.Fatalf("expected disabled=false")
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

func TestUsers_SetDisabled_notFound(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{
		users: []model.User{{ID: testAdminID, Role: "admin"}},
	})
	_, err := svc.SetDisabled(context.Background(), testAdminID, testDevID, true)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
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

func TestUsers_ResetPassword_invalidID(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{
		users: []model.User{{ID: testAdminID, Role: "admin"}},
	})
	err := svc.ResetPassword(context.Background(), testAdminID, "TARGET_USER_ID", "newpassword")
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}

func TestUsers_ResetPassword_notFound(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{
		users: []model.User{{ID: testAdminID, Role: "admin"}},
	})
	err := svc.ResetPassword(context.Background(), testAdminID, testDevID, "newpassword")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestUsers_SoftDelete_forbiddenSelf(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{
		users: []model.User{{ID: testAdminID, Role: "admin"}},
	})
	err := svc.SoftDelete(context.Background(), testAdminID, testAdminID)
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestUsers_SoftDelete_ok(t *testing.T) {
	repo := &fakeUserAdminRepo{
		users: []model.User{
			{ID: testAdminID, Role: "admin"},
			{ID: testDevID, Role: "developer"},
		},
	}
	svc := NewUsers(repo)
	err := svc.SoftDelete(context.Background(), testAdminID, testDevID)
	if err != nil {
		t.Fatal(err)
	}
	if len(repo.users) != 1 || repo.users[0].ID != testAdminID {
		t.Fatalf("expected only admin left, got %+v", repo.users)
	}
}

func TestUsers_SoftDelete_invalidID(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{
		users: []model.User{{ID: testAdminID, Role: "admin"}},
	})
	err := svc.SoftDelete(context.Background(), testAdminID, "TARGET_USER_ID")
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}
}

func TestUsers_SoftDelete_notFound(t *testing.T) {
	svc := NewUsers(&fakeUserAdminRepo{
		users: []model.User{{ID: testAdminID, Role: "admin"}},
	})
	err := svc.SoftDelete(context.Background(), testAdminID, testDevID)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
