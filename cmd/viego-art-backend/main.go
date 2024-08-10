package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ilholbea/viego-art-backend/config"
	"github.com/ilholbea/viego-art-backend/database"
	"github.com/ilholbea/viego-art-backend/handlers"
	"github.com/ilholbea/viego-art-backend/repository"
	"github.com/ilholbea/viego-art-backend/services"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewDatabase(conf.DatabaseConfig)
	if err != nil {
		log.Fatalf("Could not connect to the db: %v\n", err)
	}
	defer db.Close()

	// Init repositories
	accountRepo := repository.NewAccountRepository(db)

	// Init services
	accountService := services.NewAccountService(accountRepo)

	// Init routes
	routes := handlers.NewRoutes(accountService)

	// Init router
	router := handlers.NewRouter(routes)
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Set log level
	log.SetLevel(conf.LogLevel)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Could not listen on port ':8080': %s\n", err)
		}
	}()
	fmt.Println("Server is ready to handle requests at :8080")

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	log.Println("Server stopped")

}
