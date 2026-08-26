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

type DeploymentsHandler struct {
	deployments *service.Deployments
}

func NewDeploymentsHandler(d *service.Deployments) *DeploymentsHandler {
	return &DeploymentsHandler{deployments: d}
}

func (h *DeploymentsHandler) Create(w http.ResponseWriter, r *http.Request) {
	serviceID := r.PathValue("id")

	var req model.CreateDeploymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	out, err := h.deployments.Create(r.Context(), serviceID, req)
	if err != nil {
		writeDeploymentError(w, "Create deployment", err)
		return
	}
	writeJSON(w, http.StatusCreated, out)
}

func (h *DeploymentsHandler) List(w http.ResponseWriter, r *http.Request) {
	out, err := h.deployments.List(r.Context())
	if err != nil {
		log.Printf("List deployments: %v", err)
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *DeploymentsHandler) ListByService(w http.ResponseWriter, r *http.Request) {
	serviceID := r.PathValue("id")
	out, err := h.deployments.ListByService(r.Context(), serviceID)
	if err != nil {
		writeDeploymentError(w, "List deployments by service", err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *DeploymentsHandler) Get(w http.ResponseWriter, r *http.Request) {
	out, err := h.deployments.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeDeploymentError(w, "Get deployment", err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *DeploymentsHandler) ClaimNext(w http.ResponseWriter, r *http.Request) {
	var req struct {
		WorkerID string `json:"worker_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	job, err := h.deployments.ClaimNext(r.Context(), req.WorkerID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeNoContent(w)
			return
		}
		writeDeploymentError(w, "Claim job", err)
		return
	}
	writeJSON(w, http.StatusOK, job)
}

func (h *DeploymentsHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req model.UpdateDeploymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	out, err := h.deployments.Update(r.Context(), r.PathValue("id"), req)
	if err != nil {
		writeDeploymentError(w, "Update deployment", err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *DeploymentsHandler) Redeploy(w http.ResponseWriter, r *http.Request) {
	out, err := h.deployments.Redeploy(r.Context(), r.PathValue("id"))
	if err != nil {
		writeDeploymentError(w, "Redeploy", err)
		return
	}
	writeJSON(w, http.StatusCreated, out)
}

func writeDeploymentError(w http.ResponseWriter, op string, err error) {
	switch {
	case errors.Is(err, service.ErrInvalid):
		msg := err.Error()
		if i := strings.Index(msg, ": "); i >= 0 {
			msg = msg[i+2:]
		}
		writeError(w, http.StatusBadRequest, msg)
	case errors.Is(err, service.ErrNotFound):
		writeError(w, http.StatusNotFound, "Not found")
	case errors.Is(err, service.ErrConflict):
		msg := err.Error()
		if i := strings.Index(msg, ": "); i >= 0 {
			msg = msg[i+2:]
		}
		writeError(w, http.StatusConflict, msg)
	default:
		log.Printf("%s: %v", op, err)
		writeError(w, http.StatusInternalServerError, "Internal server error")
	}
}
