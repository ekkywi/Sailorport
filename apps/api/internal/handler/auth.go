package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
	"github.com/ekkywi/sailorport/apps/api/internal/service"
)

type AuthHandler struct {
	auth *service.Auth
}

func NewAuthHandler(auth *service.Auth) *AuthHandler {
	return &AuthHandler{auth: auth}
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

	res, err := h.auth.Login(r.Context(), req)
	if err != nil {
		writeCatalogError(w, "login", err)
		return
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
