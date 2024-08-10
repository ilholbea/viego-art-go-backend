package handlers

import (
	"context"
	"encoding/json"
	"github.com/ilholbea/viego-art-backend/types"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/ilholbea/viego-art-backend/services"
)

type Routes struct {
	accountService services.AccountService
}

func NewRoutes(accountService services.AccountService) *Routes {
	return &Routes{accountService: accountService}
}

func NewRouter(routes *Routes) http.Handler {
	router := chi.NewRouter()

	router.Get("/health", routes.Health)
	router.Get("/accounts/{email}", routes.GetAccount)
	router.Post("/accounts", routes.CreateAccount)

	return router
}

// Health handles HTTP GET requests for checking the health status of the service.
func (h *Routes) Health(w http.ResponseWriter, r *http.Request) {
	log.Debug("Health Endpoint Hit")
	w.WriteHeader(http.StatusOK)
}

func (h *Routes) GetAccount(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	log.Debug("Get Account Endpoint Hit")
	panic("implement me")
}

func (h *Routes) CreateAccount(w http.ResponseWriter, r *http.Request) {
	log.Debug("Create Account Endpoint Hit")
	var newAccount types.Account
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Decode the payload
	err := json.NewDecoder(r.Body).Decode(&newAccount)
	if err != nil {
		log.WithError(err).Error("Error decoding body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Debugf("New Account Creation Request: %v", newAccount)

	// Call the service to create the account
	if err = h.accountService.Create(ctx, newAccount); err != nil {
		log.WithError(err).Error("Error creating account")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Prepare and write the response
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Account created successfully",
	})
}
