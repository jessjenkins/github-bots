package api

import (
	"context"
	"github.com/gorilla/mux"
	"log"
)

//API provides a struct to wrap the api around
type API struct {
	Router *mux.Router
}

func Init(ctx context.Context, r *mux.Router) *API {
	api := &API{
		Router: r,
	}

	r.HandleFunc("/hello", HelloHandler()).Methods("GET")
	return api
}

func (*API) Close(ctx context.Context) error {
	// Close any dependencies
	log.Printf("graceful shutdown of api complete")
	return nil
}
