package api

import (
	"encoding/json"
	"log"
	"net/http"
)

const helloMessage = "Hello, World!"

type HelloResponse struct {
	Message string `json:"message,omitempty"`
}

// HelloHandler returns function containing a simple hello world example of an api handler
func HelloHandler(slack slackClient) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		log.Printf("Hello called")

		response := HelloResponse{
			Message: helloMessage,
		}

		err := slack.SendDirectMessage(ctx, "jessjenkins", "moo")
		if err != nil {
			log.Printf("ERROR: sending slack message failed: %v", err)
			http.Error(w, "Failed to send slack message", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Printf("ERROR: marshalling response failed: %v", err)
			http.Error(w, "Failed to marshall json response", http.StatusInternalServerError)
			return
		}

		_, err = w.Write(jsonResponse)
		if err != nil {
			log.Printf("ERROR: writing response failed: %v", err)
			http.Error(w, "Failed to write http response", http.StatusInternalServerError)
			return
		}
	}
}
