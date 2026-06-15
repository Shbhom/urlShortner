package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/shbhom/urlShortner/internal/db/postgres"
	"github.com/shbhom/urlShortner/internal/models"
	"github.com/shbhom/urlShortner/internal/pkg/metrics"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func (s *Server) AddUrlHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var body models.CreateUrlDTO

		if err := s.Services.ParseBody(r.Body, &body); err != nil {
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				slog.Error("Received Invalid request body", "error", err.Error())
				s.RespondMessage(w, &RespondMessage{Message: "Received Invalid request body: either url missing or not https"}, http.StatusBadRequest)
				return
			}
			slog.Error("Error while parsing request body: ", "error", err.Error())
			s.RespondMessage(w, &RespondMessage{Message: "Error while parsing request body"}, http.StatusBadRequest)
			return
		}
		code, err := s.Services.Addurl(r.Context(), body.Url)
		if err != nil {
			slog.Error("Error while inserting record for url", "error", err.Error())
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
		start := time.Now()
		defer func() {
			metrics.RedirectRequestHistogram.Record(r.Context(), time.Since(start).Seconds())
		}()

		code := mux.Vars(r)["code"]
		url, err := s.Services.GetUrl(r.Context(), code)
		if err != nil {
			if errors.Is(err, postgres.ErrURLNotFound) {
				metrics.RedirectNotFoundTotal.Add(r.Context(), 1)
				metrics.RedirectRequestsTotal.Add(r.Context(), 1, metric.WithAttributes(
					attribute.String("endpoint", "/r/:code"),
					attribute.String("status", "404"),
				))
				slog.Error("error while fetching target Url", "error", err.Error())
				s.RespondMessage(w, &RespondMessage{Message: "No record found for provided code"}, http.StatusNotFound)
				return
			}
			slog.Error("Error while fetching target Url", "error", err.Error())
			metrics.RedirectRequestsTotal.Add(r.Context(), 1, metric.WithAttributes(
				attribute.String("endpoint", "/r/:code"),
				attribute.String("status", "500"),
			))
			s.RespondMessage(w, &RespondMessage{Message: "Error while fetching target Url"}, http.StatusInternalServerError)
			return
		}
		metrics.RedirectionSuccessTotal.Add(r.Context(), 1)
		metrics.RedirectRequestsTotal.Add(r.Context(), 1, metric.WithAttributes(
			attribute.String("endpoint", "/r/:code"),
			attribute.String("status", "302"),
		))
		http.Redirect(w, r, url, http.StatusFound)
	}
}
