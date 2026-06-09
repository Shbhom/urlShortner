package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shbhom/urlShortner/internal/config"
	"github.com/shbhom/urlShortner/internal/services"
)

func Run(envType string) {
	config := config.LoadConfig(envType)
	r := mux.NewRouter()
	svc := services.NewService(config.DB_URL)
	serve := &Server{
		Router:   r,
		Config:   config,
		Services: svc,
	}
	serve.routes()

	switch config.API_PORT {
	case 8000:
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.API_PORT), serve))
	case 443:
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", config.API_PORT), config.ChainCertPath, config.PemCertPath, serve))
	default:
		log.Fatal("invalid API port env value found: ", config.API_PORT)
	}

}
