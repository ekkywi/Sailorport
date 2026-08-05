package main

import (
	"log"
	"net/http"

	"github.com/ekkywi/sailorport/apps/api/internal/config"
	"github.com/ekkywi/sailorport/apps/api/internal/handler"
)

func main() {
	cfg := config.Load()

	mux := http.NewServeMux()

	health := handler.NewHealthHandler("sailorport-api", cfg.Version)
	echo := handler.NewEchoHandler()

	mux.Handle("/healthz", health)
	mux.Handle("/api/v1/echo", echo)

	addr := ":" + cfg.Port
	log.Printf("Sailorport API (%s) running on http://localhost%s", cfg.AppEnv, addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}
