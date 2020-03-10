package api

import (
	"context"
	"github.com/gorilla/mux"
	"log"
)

type slackClient interface {
	SendDirectMessage(ctx context.Context, target string, message string) error
}

//API provides a struct to wrap the api around
type API struct {
	Router *mux.Router
}

func Init(slack slackClient, r *mux.Router) *API {
	api := &API{
		Router: r,
	}

	r.HandleFunc("/hello", HelloHandler(slack)).Methods("GET")
	return api
}

func (*API) Close(ctx context.Context) error {
	// Close any dependencies
	log.Printf("graceful shutdown of api complete")
	return nil
}
