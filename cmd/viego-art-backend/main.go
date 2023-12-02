package main

import (
	"github.com/ilholbea/viego-art/handlers"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	router := handlers.Router()
	log.Printf("Server Started...")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Errorf("unable to start server: %v", err)
	}
}
