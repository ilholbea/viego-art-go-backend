package handlers

import (
	"context"
	"encoding/json"
	"github.com/ilholbea/viego-art-backend/internal/services"
	"github.com/ilholbea/viego-art-backend/types"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Routes struct {
	userService services.UserService
}

func NewRoutes(userService services.UserService) *Routes {
	return &Routes{userService: userService}
}

func NewRouter(routes *Routes) http.Handler {
	router := chi.NewRouter()

	router.Get("/health", routes.Health)
	router.Post("/login", routes.Login)
	router.Post("/users", routes.CreateUser) // Signup

	return router
}

// Health handles HTTP GET requests for checking the health status of the service.
func (h *Routes) Health(w http.ResponseWriter, r *http.Request) {
	log.Debug("Health Endpoint Hit")
	w.WriteHeader(http.StatusOK)
}

func (h *Routes) Login(w http.ResponseWriter, r *http.Request) {
	log.Debug("Login Endpoint Hit")

	var receivedUser types.User
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Decode the payload
	err := json.NewDecoder(r.Body).Decode(&receivedUser)
	if err != nil {
		log.WithError(err).Error("Error decoding body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	foundUser, err := h.userService.GetUser(ctx, receivedUser)
	if err != nil {
		log.WithError(err).Error("Error getting user")
		w.WriteHeader(http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(foundUser)
	if err != nil {
		log.WithError(err).Error("Error encoding body")
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Debugf("Login User found: %v", foundUser)
}

func (h *Routes) CreateUser(w http.ResponseWriter, r *http.Request) {
	log.Debug("CreateUser User Endpoint Hit")
	var newUser types.User
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Decode the payload
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		log.WithError(err).Error("Error decoding body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Debugf("New User Creation Request: %v", newUser)

	// Call the service to create the use
	if err = h.userService.CreateUser(ctx, &newUser); err != nil {
		log.WithError(err).Error("Error creating user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Prepare and write the response
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully",
	})
}
