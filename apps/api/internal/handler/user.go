package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
	"github.com/ekkywi/sailorport/apps/api/internal/service"
)

type UsersHandler struct {
	users *service.Users
}

func NewUsersHandler(users *service.Users) *UsersHandler {
	return &UsersHandler{users: users}
}

func (h *UsersHandler) List(w http.ResponseWriter, r *http.Request) {
	out, err := h.users.List(r.Context())
	if err != nil {
		log.Printf("List users: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *UsersHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	claims := UserFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.UpdateUserRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	out, err := h.users.UpdateRole(r.Context(), claims.UserID, r.PathValue("id"), req.Role)
	if err != nil {
		writeUserError(w, "update user role", err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func writeUserError(w http.ResponseWriter, op string, err error) {
	switch {
	case errors.Is(err, service.ErrInvalid):
		msg := err.Error()
		if i := strings.Index(msg, ": "); i >= 0 {
			msg = msg[i+2:]
		}
		writeError(w, http.StatusBadRequest, msg)
	case errors.Is(err, service.ErrNotFound):
		writeError(w, http.StatusNotFound, "user not found")
	case errors.Is(err, service.ErrForbidden):
		writeError(w, http.StatusForbidden, "forbidden")
	default:
		log.Printf("%s: %v", op, err)
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
