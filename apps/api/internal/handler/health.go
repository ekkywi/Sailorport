package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type HealthHandler struct {
	Service string
	Version string
}

type healthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

func NewHealthHandler(service, version string) *HealthHandler {
	return &HealthHandler{
		Service: service,
		Version: version,
	}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := healthResponse{
		Status:  "ok",
		Service: h.Service,
		Version: h.Version,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Health encode error: %v", err)
	}
}
