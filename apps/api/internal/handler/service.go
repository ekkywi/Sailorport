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
	writeJSON(w, http.StatusOK, service.PublicServices(services))
}

func (h *ServicesHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := UserFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req model.CreateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	svc, err := h.catalog.Create(r.Context(), req, claims.UserID, claims.Email)
	if err != nil {
		writeCatalogError(w, "create service", err)
		return
	}
	writeJSON(w, http.StatusCreated, service.PublicService(svc))
}

func (h *ServicesHandler) Get(w http.ResponseWriter, r *http.Request) {
	svc, err := h.catalog.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeCatalogError(w, "get service", err)
		return
	}
	writeJSON(w, http.StatusOK, service.PublicService(svc))
}

func (h *ServicesHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := UserFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req model.UpdateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	svc, err := h.catalog.Update(r.Context(), r.PathValue("id"), req, claims.UserID, claims.Email)
	if err != nil {
		writeCatalogError(w, "update service", err)
		return
	}
	writeJSON(w, http.StatusOK, service.PublicService(svc))
}

func (h *ServicesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims := UserFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if err := h.catalog.Delete(r.Context(), r.PathValue("id"), claims.UserID, claims.Email); err != nil {
		writeCatalogError(w, "delete service", err)
		return
	}
	writeNoContent(w)
}

func writeCatalogError(w http.ResponseWriter, op string, err error) {
	switch {
	case errors.Is(err, service.ErrInvalid):
		writeError(w, http.StatusBadRequest, catalogClientMessage(err))
	case errors.Is(err, service.ErrNotFound):
		writeError(w, http.StatusNotFound, "service not found")
	case errors.Is(err, service.ErrConflict):
		writeError(w, http.StatusConflict, catalogClientMessage(err))
	case errors.Is(err, service.ErrUnauthorized):
		writeError(w, http.StatusUnauthorized, "unauthorized")
	case errors.Is(err, service.ErrForbidden):
		writeError(w, http.StatusForbidden, catalogClientMessage(err))
	default:
		log.Printf("%s: %v", op, err)
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func catalogClientMessage(err error) string {
	msg := err.Error()
	if i := strings.Index(msg, ": "); i >= 0 {
		return msg[i+2:]
	}
	return msg
}
