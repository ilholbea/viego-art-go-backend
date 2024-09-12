package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ilholbea/viego-art-backend/internal/services"
	"github.com/ilholbea/viego-art-backend/pkg/utils"
	"github.com/ilholbea/viego-art-backend/types"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

type Routes struct {
	userService services.UserService
}

func NewRoutes(userService services.UserService) *Routes {
	return &Routes{userService: userService}
}

func NewRouter(routes *Routes) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/health", routes.Health)
	router.Post("/login", routes.Login)
	router.Post("/users", routes.Signup) // Signup

	router.Group(func(r chi.Router) {
		r.Use(utils.JWTMiddleware)

		// Protected route
		router.Delete("/users", routes.DeleteUser)
	})
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

	//err = json.NewEncoder(w).Encode(foundUser)
	//if err != nil {
	//	log.WithError(err).Error("Error encoding body")
	//	w.WriteHeader(http.StatusInternalServerError)
	//}

	if foundUser != nil {
		log.Debugf("Login User found: %v", foundUser)
		tokenString, err := utils.CreateToken(foundUser.Username)
		if err != nil {
			log.WithError(err).Error("Error creating token")
			w.WriteHeader(http.StatusInternalServerError)
		}

		// Set the JWT token in a secure HTTP-only cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "jwt_token",
			Value:    tokenString,
			Expires:  time.Now().Add(1 * time.Hour),
			HttpOnly: true,                             // Prevent JavaScript access
			Secure:   os.Getenv("ENV") == "production", // Send over HTTPS only if in production
			Path:     "/",                              // Cookie is available to all endpoints
		})

		w.Write([]byte("Logged in successfully"))
		w.WriteHeader(http.StatusOK)
	} else {
		log.Debugf("Login User not found")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Invalid credentials")
	}

}

func (h *Routes) Signup(w http.ResponseWriter, r *http.Request) {
	log.Debug("Signup User Endpoint Hit")
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

func (h *Routes) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User deleted successfully"))
}
