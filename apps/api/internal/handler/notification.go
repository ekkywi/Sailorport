package handler

import (
	"net/http"
	"strconv"

	"github.com/ekkywi/sailorport/apps/api/internal/model"
	"github.com/ekkywi/sailorport/apps/api/internal/service"
)

type NotificationsHandler struct {
	deployments *service.Deployments
}

func NewNotificationsHandler(d *service.Deployments) *NotificationsHandler {
	return &NotificationsHandler{deployments: d}
}

func (h *NotificationsHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := UserFromContext(r.Context())
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	limit := 50
	if raw := r.URL.Query().Get("limit"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 {
			writeError(w, http.StatusBadRequest, "limit must be positive integer")
			return
		}
		limit = n
	}

	out, err := h.deployments.ListNotifications(r.Context(), claims.UserID, claims.Role, limit)

	if err != nil {
		writeDeploymentError(w, "List notifications", err)
		return
	}
	if out == nil {
		out = []model.Notification{}
	}
	writeJSON(w, http.StatusOK, out)
}
