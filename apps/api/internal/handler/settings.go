package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
	"github.com/ekkywi/sailorport/apps/api/internal/service"
)

type SettingsHandler struct {
	settings *service.Settings
}

func NewSettingsHandler(settings *service.Settings) *SettingsHandler {
	return &SettingsHandler{settings: settings}
}

func (h *SettingsHandler) Get(w http.ResponseWriter, r *http.Request) {
	out, err := h.settings.Get(r.Context())
	if err != nil {
		writeCatalogError(w, "get settings", err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *SettingsHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req model.UpdateAppSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	out, err := h.settings.Update(r.Context(), req)
	if err != nil {
		writeCatalogError(w, "update settings", err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *SettingsHandler) RegistrationStatus(w http.ResponseWriter, r *http.Request) {
	out, err := h.settings.RegistrationStatus(r.Context())
	if err != nil {
		writeCatalogError(w, "registration status", err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}
