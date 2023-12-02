package main

import (
	"github.com/ilholbea/viego-art/handlers"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file ")
	}

	router := handlers.Router()
	log.Printf("Server Started...")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Errorf("unable to start server: %v", err)
	}
}
