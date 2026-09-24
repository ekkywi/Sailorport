package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
	"github.com/ekkywi/sailorport/apps/api/internal/ratelimit"
	"github.com/ekkywi/sailorport/apps/api/internal/service"
	"github.com/ekkywi/sailorport/apps/api/internal/store"
)

type fakeLoginUsers struct{}

func (f *fakeLoginUsers) Count(ctx context.Context) (int, error) { return 1, nil }

func (f *fakeLoginUsers) Create(ctx context.Context, email, name, passwordHash, role string) (model.User, error) {
	return model.User{}, store.ErrConflict
}

func (f *fakeLoginUsers) GetByEmail(ctx context.Context, email string) (model.UserRecord, error) {
	return model.UserRecord{}, store.ErrNotFound
}

func (f *fakeLoginUsers) GetByID(ctx context.Context, id string) (model.User, error) {
	return model.User{}, store.ErrNotFound
}

type fakeLoginSettings struct{}

func (f *fakeLoginSettings) RegistrationOpen(ctx context.Context) (bool, error) {
	return false, nil
}

func newLoginTestHandler(limit int) *AuthHandler {
	auth := service.NewAuth(&fakeLoginUsers{}, &fakeLoginSettings{}, "test-jwt-secret-for-ratelimit")
	lim := ratelimit.New(limit, time.Minute)
	return NewAuthHandler(auth, lim)
}

func postLogin(t *testing.T, h *AuthHandler, remoteAddr string) (int, map[string]string) {
	t.Helper()
	body := `{"email":"nobody@example.com","password":"wrong"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = remoteAddr
	rr := httptest.NewRecorder()
	h.Login(rr, req)

	var out map[string]string
	_ = json.NewDecoder(rr.Body).Decode(&out)
	return rr.Code, out
}

func TestLogin_RateLimitedAfterFailures(t *testing.T) {
	h := newLoginTestHandler(3)
	addr := "203.0.113.10:54321"

	for i := 1; i <= 3; i++ {
		code, _ := postLogin(t, h, addr)
		if code != http.StatusUnauthorized {
			t.Fatalf("attempt %d: want 401, got %d", i, code)
		}
	}

	code, out := postLogin(t, h, addr)
	if code != http.StatusTooManyRequests {
		t.Fatalf("4th attempt: want 429, got %d body=%v", code, out)
	}
	if out["error"] != "too many login attempts" {
		t.Fatalf("error message: %v", out)
	}
}

func TestLogin_RateLimitIsPerIp(t *testing.T) {
	h := newLoginTestHandler(2)
	a := "203.0.113.10:1111"
	b := "203.0.113.11:2222"

	for i := 0; i < 2; i++ {
		code, _ := postLogin(t, h, a)
		if code != http.StatusUnauthorized {
			t.Fatalf("ip A attempt %d: want 401, got %d", i+1, code)
		}
	}
	code, _ := postLogin(t, h, a)
	if code != http.StatusTooManyRequests {
		t.Fatalf("ip A should be 429, got %d", code)
	}

	code, _ = postLogin(t, h, b)
	if code != http.StatusUnauthorized {
		t.Fatalf("ip B should still 401, got %d", code)
	}
}
