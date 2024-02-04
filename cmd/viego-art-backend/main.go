package main

import (
	"github.com/ilholbea/viego-art/handlers"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	server := handlers.NewServer()
	log.Printf("Server Started...")
	err := http.ListenAndServe(":8080", server.Router)
	if err != nil {
		log.Fatalf("unable to start server: %v", err)
	}
}
