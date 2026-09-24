package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"fmt"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
	"github.com/ekkywi/sailorport/apps/api/internal/ratelimit"
	"github.com/ekkywi/sailorport/apps/api/internal/service"
)

type AuthHandler struct {
	auth    *service.Auth
	logins	*ratelimit.Limiter
}

func NewAuthHandler(auth *service.Auth, logins *ratelimit.Limiter) *AuthHandler {
	return &AuthHandler{auth: auth, logins: logins}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	user, err := h.auth.Register(r.Context(), req)
	if err != nil {
		writeCatalogError(w, "register", err)
		return
	}
	writeJSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	key := clientIP(r)
	if h.logins != nil && !h.logins.Allow(key) {
		writeCatalogError(w, "login", fmt.Errorf("%w: too many login attempts", service.ErrRateLimited))
		return
	}

	res, err := h.auth.Login(r.Context(), req)
	if err != nil {
		if h.logins != nil && errors.Is(err, service.ErrUnauthorized) {
			h.logins.Fail(key)
		}
		writeCatalogError(w, "login", err)
		return
	}

	if h.logins != nil {
		h.logins.Reset(key)
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims := UserFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	user, err := h.auth.Me(r.Context(), claims.UserID)
	if err != nil {
		writeCatalogError(w, "me", err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}
