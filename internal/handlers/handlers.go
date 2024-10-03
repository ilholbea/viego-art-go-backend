package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ilholbea/viego-art-backend/internal/services"
	"github.com/ilholbea/viego-art-backend/pkg/utils"
	"github.com/ilholbea/viego-art-backend/types"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"time"
)

type Routes struct {
	userService  services.UserService
	imageService services.ImageService
}

func NewRoutes(userService services.UserService, imageService services.ImageService) *Routes {
	return &Routes{userService: userService, imageService: imageService}
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

		// Protected routes
		r.Delete("/users", routes.DeleteUser)
		r.Post("/images", routes.UploadImage)

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

	if foundUser != nil {
		log.Debugf("Login User found: %v", foundUser)
		tokenString, err := utils.CreateToken(foundUser)
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
	// TODO ADD JWT CHECK
	w.Write([]byte("User deleted successfully"))
}

func (h *Routes) UploadImage(w http.ResponseWriter, r *http.Request) {
	// TODO ADD JWT CHECK, DEBUGGER LOGS
	var newFile types.UploadedFile
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := r.ParseMultipartForm(32 << 50)
	if err != nil {
		log.WithError(err).Error("Error parsing multipart form")
		w.WriteHeader(http.StatusBadRequest)
	}

	filename := r.FormValue("filename")
	bucket := r.FormValue("bucket")

	if filename == "" {
		log.WithError(err).Error("Error parsing file name")
		w.WriteHeader(http.StatusBadRequest)
	}
	newFile.Filename = filename

	if bucket == "" {
		log.WithError(err).Error("Error parsing bucket name")
		w.WriteHeader(http.StatusBadRequest)
	}

	// Parse and handle the file
	file, _, err := r.FormFile("file")
	if err != nil {
		log.WithError(err).Error("Error retrieving file")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the file to get content type
	buff := make([]byte, 512)
	file.Read(buff)
	file.Seek(0, 0)
	contentType := http.DetectContentType(buff)
	newFile.ContentType = contentType

	// Create a buffer to store file content
	var fileBuffer bytes.Buffer
	fileSize, err := io.Copy(&fileBuffer, file)
	if err != nil {
		log.WithError(err).Error("Error reading file content")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newFile.Size = fileSize
	newFile.File = fileBuffer.String()

	err = h.imageService.UploadImage(ctx, bucket, &newFile)
	if err != nil {
		log.WithError(err).Error("Error uploading image")
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
