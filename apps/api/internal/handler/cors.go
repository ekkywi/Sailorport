package handler

import (
	"net/http"
	"slices"
)

const corsAllowMethods = "GET, POST, PUT, PATCH, DELETE, OPTIONS"
const corsAllowHeaders = "Content-Type, Authorization"

func CORS(next http.Handler, allowedOrigins []string) http.Handler {
	allowed := slices.Clone(allowedOrigins)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && slices.Contains(allowed, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}

		w.Header().Set("Access-Control-Allow-Methods", corsAllowMethods)
		w.Header().Set("Access-Control-Allow-Headers", corsAllowHeaders)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
