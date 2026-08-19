package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ekkywi/sailorport/apps/api/internal/service"
)

type ScaffoldHandler struct {
	scaffold *service.Scaffold
}

func NewScaffoldHandler(scaffold *service.Scaffold) *ScaffoldHandler {
	return &ScaffoldHandler{scaffold: scaffold}
}

func (h *ScaffoldHandler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	items, err := h.scaffold.ListTemplates()
	if err != nil {
		log.Printf("list templates: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *ScaffoldHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := UserFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req service.ScaffoldRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	result, err := h.scaffold.Run(r.Context(), req, claims.UserID, claims.Email)
	if err != nil {
		writeCatalogError(w, "scaffold service", err)
		return
	}
	writeJSON(w, http.StatusCreated, result)
}
