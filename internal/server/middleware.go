package server

import (
	"log/slog"
	"net/http"
	"strings"
)

func (s *Server) enableCors(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Running in enableCors middleware")

		origin := r.Header.Get("Origin")
		allowedOrigins := []string{
			"http://localhost:8000",
			getBaseURL(r),
		}

		originAllowed := false
		// Allow localhost with any port in development
		if strings.HasPrefix(origin, "http://localhost:") || origin == "http://localhost" || origin == "" {
			originAllowed = true
		}

		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				originAllowed = true
				break
			}
		}

		if !originAllowed {
			http.Error(w, "Origin not allowed", http.StatusForbidden)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin, X-Api-Key, X-Requested-With, Accept")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler(w, r)
	}
}
