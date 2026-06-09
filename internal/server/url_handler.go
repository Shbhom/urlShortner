package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/shbhom/urlShortner/internal/models"
)

func (s *Server) AddUrlHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body models.CreateUrlDTO

		if err := s.Services.ParseBody(r.Body, &body); err != nil {
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				slog.Error("Received Invalid request body", err.Error(), "error")
				s.RespondMessage(w, &RespondMessage{Message: "Received Invalid request body: either url missing or not https"}, http.StatusBadRequest)
				return
			}
			slog.Error("Error while parsing request body: ", err.Error(), "error")
			s.RespondMessage(w, &RespondMessage{Message: "Error while parsing request body"}, http.StatusBadRequest)
			return
		}
		code, err := s.Services.Addurl(body.Url)
		if err != nil {
			slog.Error("Error while inserting record for url", err.Error(), "error")
			s.RespondMessage(w, &RespondMessage{Message: "Error while adding record to db"}, http.StatusInternalServerError)
			return
		}
		shortUrl := fmt.Sprintf("%s/r/%s", s.Config.BASE_URL, code)
		var resp struct {
			ShortUrl string `json:"shortUrl"`
		}
		resp.ShortUrl = shortUrl
		s.RespondMessage(w, &RespondMessage{Message: "Successfull", Data: resp}, http.StatusCreated)
	}
}

func (s *Server) RedirectionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("handling redirection request")
		code := mux.Vars(r)["code"]
		url, err := s.Services.GetUrl(code)
		if err != nil {
			slog.Error("Error while fetching target Url", err.Error(), "error")
			s.RespondMessage(w, &RespondMessage{Message: "Error while fetching target Url"}, http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, url.String(), http.StatusPermanentRedirect)
	}
}
