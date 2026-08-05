package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type echoRequest struct {
	Message string `json:"message"`
}

type echoResponse struct {
	Reply string `json:"reply"`
}

type EchoHandler struct {}

func NewEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (h *EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req echoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid json body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Message == "" {
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}

	response := echoResponse{
		Reply: "Sailorport received: " + req.Message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Echo encode error: %v", err)
	}
}