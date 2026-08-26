package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/ekkywi/sailorport/apps/api/internal/service"
)

type CatalogAppsHandler struct {
	apps *service.CatalogApps
}

func NewCatalogAppsHandler(apps *service.CatalogApps) *CatalogAppsHandler {
	return &CatalogAppsHandler{apps: apps}
}

func (h *CatalogAppsHandler) List(w http.ResponseWriter, r *http.Request) {
	out, err := h.apps.List(r.Context())
	if err != nil {
		log.Printf("List catalog apps: %v", err)
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *CatalogAppsHandler) Get(w http.ResponseWriter, r *http.Request) {
	out, err := h.apps.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Not found")
			return
		}
		log.Printf("Get catalog app: %v", err)
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, out)
}
