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

type ServicesHandler struct {
	catalog *service.Catalog
}

func NewServicesHandler(catalog *service.Catalog) *ServicesHandler {
	return &ServicesHandler{catalog: catalog}
}

func (h *ServicesHandler) List(w http.ResponseWriter, r *http.Request) {
	services, err := h.catalog.List(r.Context())
	if err != nil {
		log.Printf("list services: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, services)
}

func (h *ServicesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	svc, err := h.catalog.Create(r.Context(), req)
	if err != nil {
		writeCatalogError(w, "create service", err)
		return
	}
	writeJSON(w, http.StatusCreated, svc)
}

func (h *ServicesHandler) Get(w http.ResponseWriter, r *http.Request) {
	svc, err := h.catalog.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeCatalogError(w, "get service", err)
		return
	}
	writeJSON(w, http.StatusOK, svc)
}

func (h *ServicesHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req model.UpdateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	svc, err := h.catalog.Update(r.Context(), r.PathValue("id"), req)
	if err != nil {
		writeCatalogError(w, "update service", err)
		return
	}
	writeJSON(w, http.StatusOK, svc)
}

func (h *ServicesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.catalog.Delete(r.Context(), r.PathValue("id")); err != nil {
		writeCatalogError(w, "delete service", err)
		return
	}
	writeNoContent(w)
}

func writeCatalogError(w http.ResponseWriter, op string, err error) {
	switch {
	case errors.Is(err, service.ErrInvalid):
		msg := err.Error()
		if i := strings.Index(msg, ": "); i >= 0 {
			msg = msg[i+2:]
		}
		writeError(w, http.StatusBadRequest, msg)
	case errors.Is(err, service.ErrNotFound):
		writeError(w, http.StatusNotFound, "service not found")
	case errors.Is(err, service.ErrConflict):
		writeError(w, http.StatusConflict, "service already exists")
	default:
		log.Printf("%s: %v", op, err)
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
