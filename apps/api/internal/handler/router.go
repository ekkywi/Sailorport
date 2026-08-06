package handler

import (
	"net/http"

	"github.com/ekkywi/sailorport/apps/api/internal/service"
)

type API struct {
	Version  string
	Catalog  *service.Catalog
	Scaffold *service.Scaffold
}

func NewRouter(api API) http.Handler {
	mux := http.NewServeMux()

	health := NewHealthHandler("sailorport-api", api.Version)
	echo := NewEchoHandler() // demo endpoint; remove later
	services := NewServicesHandler(api.Catalog)
	scaffold := NewScaffoldHandler(api.Scaffold)

	mux.Handle("/healthz", health)
	mux.Handle("/api/v1/echo", echo)
	mux.HandleFunc("GET /api/v1/services", services.List)
	mux.HandleFunc("POST /api/v1/services", services.Create)
	mux.HandleFunc("GET /api/v1/services/{id}", services.Get)
	mux.HandleFunc("PUT /api/v1/services/{id}", services.Update)
	mux.HandleFunc("DELETE /api/v1/services/{id}", services.Delete)
	mux.HandleFunc("GET /api/v1/templates", scaffold.ListTemplates)
	mux.HandleFunc("POST /api/v1/scaffold", scaffold.Create)

	return CORS(mux)
}
