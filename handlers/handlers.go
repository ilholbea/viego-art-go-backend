package handlers

import (
	"net/http"

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

	return router
}

// Health handles HTTP GET requests for checking the health status of the service.
func (h *Routes) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Routes) GetAccount(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
