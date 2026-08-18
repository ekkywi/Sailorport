package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
	"github.com/ekkywi/sailorport/apps/api/internal/service"
)

type RuntimeHandler struct {
	runtime *service.Runtime
}

func NewRuntimeHandler(r *service.Runtime) *RuntimeHandler {
	return &RuntimeHandler{runtime: r}
}

func (h *RuntimeHandler) Stop(w http.ResponseWriter, r *http.Request) {
	var req model.RuntimeActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	job, err := h.runtime.RequestStop(r.Context(), r.PathValue("id"), req.Environment)
	if err != nil {
		writeRuntimeError(w, "Stop service", err)
		return
	}
	writeJSON(w, http.StatusAccepted, job)
}

func (h *RuntimeHandler) Start(w http.ResponseWriter, r *http.Request) {
	var req model.RuntimeActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	job, err := h.runtime.RequestStart(r.Context(), r.PathValue("id"), req.Environment)
	if err != nil {
		writeRuntimeError(w, "Start service", err)
		return
	}
	writeJSON(w, http.StatusAccepted, job)
}

func (h *RuntimeHandler) ClaimNext(w http.ResponseWriter, r *http.Request) {
	var req struct {
		WorkerID string `json:"worker_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	job, err := h.runtime.ClaimNext(r.Context(), req.WorkerID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeNoContent(w)
			return
		}
		writeRuntimeError(w, "Claim runtime job", err)
		return
	}
	writeJSON(w, http.StatusOK, job)
}

func (h *RuntimeHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req model.UpdateRuntimeJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	out, err := h.runtime.UpdateFromAgent(r.Context(), r.PathValue("id"), req)
	if err != nil {
		writeRuntimeError(w, "Update runtime job", err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func writeRuntimeError(w http.ResponseWriter, op string, err error) {
	switch {
	case errors.Is(err, service.ErrInvalid):
		msg := err.Error()
		if i := strings.Index(msg, ": "); i >= 0 {
			msg = msg[i+2:]
		}
		writeError(w, http.StatusBadRequest, msg)
	case errors.Is(err, service.ErrNotFound):
		writeError(w, http.StatusNotFound, "Not found")
	default:
		log.Printf("%s: %v", op, err)
		writeError(w, http.StatusInternalServerError, "Internal server error")
	}
}
