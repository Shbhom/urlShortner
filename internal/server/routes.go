package server

import (
	"encoding/json"
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
