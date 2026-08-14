package handler

import (
	"log"
	"net/http"

	"github.com/ekkywi/sailorport/apps/api/internal/service"
)

type EnvironmentsHandler struct {
	envs *service.Environments
}

func NewEnvironmentsHandler(envs *service.Environments) *EnvironmentsHandler {
	return &EnvironmentsHandler{envs: envs}
}

func (h *EnvironmentsHandler) List(w http.ResponseWriter, r *http.Request) {
	out, err := h.envs.List(r.Context())
	if err != nil {
		log.Printf("List environments: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, out)
}
