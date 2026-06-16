package server

import (
	"encoding/json"
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shbhom/urlShortner/internal/config"
	"github.com/shbhom/urlShortner/internal/services"
)

type Server struct {
	Router   *mux.Router
	Config   *config.Config
	Services *services.Service
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.Router.HandleFunc("/url", s.enableCors(s.AddUrlHandler())).Methods(http.MethodPost, http.MethodOptions)
	s.Router.HandleFunc("/r/{code}", s.RedirectionHandler()).Methods(http.MethodGet)

	// Extract the subdirectory from embedded files
	distFS, err := fs.Sub(embeddedUI, "ui/dist")
	if err != nil {
		// Log error but allow API routes to function if UI is missing locally
		slog.Error("Failed to extract embedded UI filesystem", "error", err)
	} else {
		spa := newSPAHandler(distFS, "index.html")
		s.Router.PathPrefix("/").Handler(spa)
	}
}

type RespondMessage struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (s *Server) RespondMessage(w http.ResponseWriter, response *RespondMessage, code int) {
	w.WriteHeader(code)
	if response != nil {
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			slog.Error("Error in encoding the response", "error", err)
			return
		}
	}
}
