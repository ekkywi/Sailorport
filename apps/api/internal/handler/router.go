package handler

import (
	"net/http"

	"github.com/ekkywi/sailorport/apps/api/internal/service"
)

type API struct {
	Version   string
	JWTSecret string
	Catalog   *service.Catalog
	Scaffold  *service.Scaffold
	Auth      *service.Auth
}

func NewRouter(api API) http.Handler {
	mux := http.NewServeMux()
	secret := api.JWTSecret

	health := NewHealthHandler("sailorport-api", api.Version)
	authH := NewAuthHandler(api.Auth)
	services := NewServicesHandler(api.Catalog)
	scaffold := NewScaffoldHandler(api.Scaffold)

	writer := []string{"developer", "admin"}
	reader := []string{"viewer", "developer", "admin"}

	mux.Handle("/healthz", health)

	mux.HandleFunc("POST /api/v1/auth/register", authH.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authH.Login)
	mux.Handle("GET /api/v1/auth/me", withAuth(secret, authH.Me))

	mux.Handle("GET /api/v1/services", withRole(secret, reader, services.List))
	mux.Handle("POST /api/v1/services", withRole(secret, writer, services.Create))
	mux.Handle("GET /api/v1/services/{id}", withRole(secret, reader, services.Get))
	mux.Handle("PUT /api/v1/services/{id}", withRole(secret, writer, services.Update))
	mux.Handle("DELETE /api/v1/services/{id}", withRole(secret, writer, services.Delete))

	mux.Handle("GET /api/v1/templates", withRole(secret, reader, scaffold.ListTemplates))
	mux.Handle("POST /api/v1/scaffold", withRole(secret, writer, scaffold.Create))

	return CORS(mux)
}
