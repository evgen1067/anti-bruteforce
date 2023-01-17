package server

import (
	"github.com/evgen1067/anti-bruteforce/internal/config"
	"github.com/evgen1067/anti-bruteforce/internal/service"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"time"
)

var s *service.Services

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/auth", Auth).Methods(http.MethodPost)
	router.HandleFunc("/list/{value}", Add).Methods(http.MethodPost)
	router.HandleFunc("/list/{value}", Delete).Methods(http.MethodDelete)

	router.Use(headersMiddleware)

	router.NotFoundHandler = router.NewRoute().HandlerFunc(CustomNotFoundHandler).GetHandler()

	return router
}

func InitServer(_s *service.Services, cfg *config.Config) *http.Server {
	s = _s
	return &http.Server{
		Addr:              net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port),
		Handler:           Router(),
		ReadHeaderTimeout: 3 * time.Second,
	}
}
