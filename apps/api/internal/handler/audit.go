package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ekkywi/sailorport/apps/api/internal/service"
)

type AuditHandler struct {
	audit *service.Audit
}

func NewAuditHandler(audit *service.Audit) *AuditHandler {
	return &AuditHandler{audit: audit}
}

func (h *AuditHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if raw := r.URL.Query().Get("limit"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 {
			writeError(w, http.StatusBadRequest, "limit must be a positive integer")
			return
		}
		limit = n
	}

	out, err := h.audit.List(r.Context(), limit)
	if err != nil {
		log.Printf("List audit: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, out)
}
