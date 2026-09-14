package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
	"github.com/ekkywi/sailorport/apps/api/internal/service"
)

type SetupHandler struct {
	setup *service.Setup
}

func NewSetupHandler(setup *service.Setup) *SetupHandler {
	return &SetupHandler{setup: setup}
}

func (h *SetupHandler) Status(w http.ResponseWriter, r *http.Request) {
	st, err := h.setup.Status(r.Context())
	if err != nil {
		writeCatalogError(w, "setup status", err)
		return
	}
	writeJSON(w, http.StatusOK, st)
}

func (h *SetupHandler) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	var req model.SetupAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	user, err := h.setup.CreateAdmin(r.Context(), req)
	if err != nil {
		writeCatalogError(w, "setup admin", err)
		return
	}
	writeJSON(w, http.StatusCreated, user)
}