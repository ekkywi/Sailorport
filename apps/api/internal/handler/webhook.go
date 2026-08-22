package handler

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/ekkywi/sailorport/apps/api/internal/service"
)

const maxWebhookBody = 1 << 20 // 1 MiB

type WebhookHandler struct {
	webhooks *service.Webhook
}

func NewWebhookHandler(w *service.Webhook) *WebhookHandler {
	return &WebhookHandler{webhooks: w}
}

func (h *WebhookHandler) GitHub(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(io.LimitReader(r.Body, maxWebhookBody+1))
	if err != nil {
		writeError(w, http.StatusBadRequest, "could not read body")
		return
	}
	if len(body) > maxWebhookBody {
		writeError(w, http.StatusRequestEntityTooLarge, "body too large")
		return
	}

	event := r.Header.Get("X-GitHub-Event")
	sig := r.Header.Get("X-Hub-Signature-256")

	ack, err := h.webhooks.HandleGitHub(r.Context(), event, sig, body)
	if err != nil {
		writeWebhookError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ack)
}

func writeWebhookError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalid):
		msg := err.Error()
		if i := strings.Index(msg, ": "); i >= 0 {
			msg = msg[i+2:]
		}
		writeError(w, http.StatusBadRequest, msg)
	case errors.Is(err, service.ErrUnauthorized):
		writeError(w, http.StatusUnauthorized, "unauthorized")
	default:
		log.Printf("webhook: %v", err)
		writeError(w, http.StatusInternalServerError, "Internal server error")
	}
}
