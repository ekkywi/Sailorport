package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
	"github.com/ekkywi/sailorport/apps/api/internal/store"
)

type ServicesHandler struct {
	store *store.ServicesStore
}

func NewServicesHandler(serviceStore *store.ServicesStore) *ServicesHandler {
	return &ServicesHandler{store: serviceStore}
}

func (h *ServicesHandler) List(w http.ResponseWriter, r *http.Request) {
	services, err := h.store.List(r.Context())
	if err != nil {
		log.Printf("List service: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, services)
}

func (h *ServicesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	req.Owner = strings.TrimSpace(req.Owner)

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	svc, err := h.store.Create(r.Context(), req)
	if errors.Is(err, store.ErrConflict) {
		http.Error(w, "Service already exists", http.StatusConflict)
		return
	}
	if err != nil {
		log.Printf("Create service: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, svc)
}

func (h *ServicesHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	svc, err := h.store.Get(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		http.Error(w, "service not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("get service: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, svc)
}
func (h *ServicesHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	var req model.UpdateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	req.Owner = strings.TrimSpace(req.Owner)
	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	svc, err := h.store.Update(r.Context(), id, req)
	if errors.Is(err, store.ErrNotFound) {
		http.Error(w, "service not found", http.StatusNotFound)
		return
	}
	if errors.Is(err, store.ErrConflict) {
		http.Error(w, "service name already exists", http.StatusConflict)
		return
	}
	if err != nil {
		log.Printf("update service: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, svc)
}
func (h *ServicesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	err := h.store.Delete(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		http.Error(w, "service not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("delete service: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("encode json: %v", err)
	}
}