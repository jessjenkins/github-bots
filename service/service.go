package service

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/jessjenkins/github-bots/api"
	"github.com/jessjenkins/github-bots/config"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"
)

var (
	gracefulShutdownTimeout = 5 * time.Second
)

type Service struct {
	Config *config.Config
	Server *http.Server
	Router *mux.Router
	API    *api.API
}

func Create() (*Service, error) {
	log.Println("creating service")
	ctx := context.Background()

	cfg, err := config.Get()
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve service configuration")
	}
	log.Println("got service configuration", cfg)

	r := mux.NewRouter()

	s := &http.Server{
		Handler:           r,
		Addr:              cfg.BindAddr,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
	}

	a := api.Init(ctx, r)

	return &Service{
		Config: cfg,
		Router: r,
		API:    a,
		Server: s,
	}, nil
}

// Run the service
func (svc *Service) Run(svcErrors chan error) error {
	log.Println("running service")

	go func() {
		if err := svc.Server.ListenAndServe(); err != nil {
			svcErrors <- errors.Wrap(err, "failure in http listen and serve")
		}
	}()

	return nil
}

// Gracefully shutdown the service
func (svc *Service) Close() {
	timeout := gracefulShutdownTimeout
	log.Printf("commencing graceful shutdown: timeout[%v]", timeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// stop any incoming requests before closing any outbound connections
	if err := svc.Server.Shutdown(ctx); err != nil {
		log.Printf("ERROR: failed to shutdown http server: %v", err)
	}

	if err := svc.API.Close(ctx); err != nil {
		log.Printf("ERROR: error closing API: %v", err)
	}

	log.Println("graceful shutdown complete")
}
