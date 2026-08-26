package handler

import (
	"context"
	"crypto/subtle"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/auth"
	"github.com/ekkywi/sailorport/apps/api/internal/model"
	"github.com/ekkywi/sailorport/apps/api/internal/service"
)

type ctxKey int

const claimsKey ctxKey = 1

// currentUserLookup mengambil user segar dari DB; dipenuhi oleh *service.Auth.
type currentUserLookup interface {
	Me(ctx context.Context, userID string) (model.User, error)
}

func UserFromContext(ctx context.Context) *auth.Claims {
	c, _ := ctx.Value(claimsKey).(*auth.Claims)
	return c
}

// RequireAuth memvalidasi JWT lalu memuat ulang user dari DB. Token berumur 24 jam,
// jadi role/disabled di dalam klaim bisa sudah basi: tanpa pembacaan ini, user yang
// di-disable, di-soft-delete, atau diturunkan role-nya tetap punya akses sampai
// tokennya kedaluwarsa.
func RequireAuth(jwtSecret string, users currentUserLookup) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				writeError(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			raw := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
			claims, err := auth.ParseToken(raw, jwtSecret)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			user, err := users.Me(r.Context(), claims.UserID)
			if err != nil {
				switch {
				case errors.Is(err, service.ErrUnauthorized),
					errors.Is(err, service.ErrNotFound),
					errors.Is(err, service.ErrInvalid):
					writeError(w, http.StatusUnauthorized, "unauthorized")
				default:
					log.Printf("load current user: %v", err)
					writeError(w, http.StatusInternalServerError, "internal server error")
				}
				return
			}

			// DB yang berwenang, bukan isi token.
			claims.Role = user.Role
			claims.Email = user.Email

			ctx := context.WithValue(r.Context(), claimsKey, &claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := map[string]struct{}{}
	for _, role := range roles {
		allowed[role] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := UserFromContext(r.Context())
			if claims == nil {
				writeError(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			if _, ok := allowed[claims.Role]; !ok {
				writeError(w, http.StatusForbidden, "forbidden")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func withAuth(secret string, users currentUserLookup, h http.HandlerFunc) http.Handler {
	return RequireAuth(secret, users)(h)
}

func withRole(secret string, users currentUserLookup, roles []string, h http.HandlerFunc) http.Handler {
	return RequireAuth(secret, users)(RequireRole(roles...)(h))
}

func withAgentToken(expected string, h http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if expected == "" {
			writeError(w, http.StatusUnauthorized, "agent token not configured")
			return
		}
		header := r.Header.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		got := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		if len(got) != len(expected) ||
			subtle.ConstantTimeCompare([]byte(got), []byte(expected)) != 1 {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		h.ServeHTTP(w, r)
	})
}
