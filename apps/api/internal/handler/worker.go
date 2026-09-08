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

type WorkersHandler struct {
	workers *service.Workers
}

func NewWorkersHandler(w *service.Workers) *WorkersHandler {
	return &WorkersHandler{workers: w}
}

func (h *WorkersHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterWorkerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	out, err := h.workers.Register(r.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalid) {
			writeError(w, http.StatusBadRequest, "Name is required")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to register worker")
		return
	}
	writeJSON(w, http.StatusCreated, out)
}

func (h *WorkersHandler) Heartbeat(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req model.HeartbeatRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	out, err := h.workers.Heartbeat(r.Context(), id, req.Status)
	if err != nil {
		if errors.Is(err, service.ErrInvalid) {
			writeError(w, http.StatusBadRequest, "Invalid heartbeat")
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Worker not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to heartbeat")
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *WorkersHandler) List(w http.ResponseWriter, r *http.Request) {
	out, err := h.workers.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list workers")
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func writeWorkerError(w http.ResponseWriter, op string, err error) {
	switch {
	case errors.Is(err, service.ErrInvalid):
		msg := err.Error()
		if i := strings.Index(msg, ": "); i >= 0 {
			msg = msg[i+2:]
		}
		writeError(w, http.StatusBadRequest, msg)
	case errors.Is(err, service.ErrNotFound):
		writeError(w, http.StatusNotFound, "Worker not found")
	default:
		log.Printf("%s: %v", op, err)
		writeError(w, http.StatusInternalServerError, "Internal server error")
	}
}

func (h *WorkersHandler) UpdateLabels(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req model.UpdateWorkerLabelsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	defer r.Body.Close()

	out, err := h.workers.UpdateLabels(r.Context(), id, req)
	if err != nil {
		writeWorkerError(w, "update worker labels", err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *WorkersHandler) Decommission(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	out, err := h.workers.Decommission(r.Context(), id)
	if err != nil {
		writeWorkerError(w, "decommission worker", err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *WorkersHandler) Restore(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	out, err := h.workers.Restore(r.Context(), id)
	if err != nil {
		writeWorkerError(w, "restore worker", err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}
