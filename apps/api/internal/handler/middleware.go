package handler

import (
	"context"
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/auth"
)

type ctxKey int

const claimsKey ctxKey = 1

func UserFromContext(ctx context.Context) *auth.Claims {
	c, _ := ctx.Value(claimsKey).(*auth.Claims)
	return c
}

func RequireAuth(jwtSecret string) func(http.Handler) http.Handler {
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

func withAuth(secret string, h http.HandlerFunc) http.Handler {
	return RequireAuth(secret)(h)
}

func withRole(secret string, roles []string, h http.HandlerFunc) http.Handler {
	return RequireAuth(secret)(RequireRole(roles...)(h))
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
