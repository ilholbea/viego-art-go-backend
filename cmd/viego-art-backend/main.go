package main

import (
	"github.com/ilholbea/viego-art/handlers"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	server, err := handlers.NewServer()
	if err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
	log.Printf("Server Started...")
	err = http.ListenAndServe(":8080", server.Router)
	if err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
}
