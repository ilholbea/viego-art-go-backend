package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/ilholbea/viego-art/config"
	"github.com/ilholbea/viego-art/internal"
	"github.com/ilholbea/viego-art/types"
)

// JwtExpirationTime Used to set the expiration time of a JWT in minutes.
const JwtExpirationTime = 1

var jwtKey = []byte(LoadEnvironmentVariable("JWT_SECRET_KEY"))
var dbConfig = &config.Database{
	Host:     LoadEnvironmentVariable("HOST"),
	Port:     LoadEnvironmentVariable("PORT"),
	User:     LoadEnvironmentVariable("USER"),
	Password: LoadEnvironmentVariable("PASSWORD"),
	Database: LoadEnvironmentVariable("DATABASE"),
	Schema:   LoadEnvironmentVariable("SCHEMA"),
}

type Credentials struct {
	Email          string `json:"email"`
	HashedPassword string `json:"hashedPassword"`
}

type CookieContent struct {
	Token          string    `json:"token"`
	ExpirationTime time.Time `json:"expirationTime"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Router() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/health", health)
	router.Post("/login", login)
	router.Route("/account", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Post("/", registerAccount)
		r.Delete("/{id}", removeAccount)
		r.Put("/", updateAccount)
	})

	return router
}

// authMiddleware Used to verify if the request has a valid JWT attached.
func authMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		err := validateJWT(request)
		if err != nil {
			log.Errorf("authMiddleware - error in authentication: %v", err)
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		handler.ServeHTTP(writer, request)
	})
}

// health Used to check the health status of the application.
func health(w http.ResponseWriter, _ *http.Request) {
	log.Printf("health - endpoint was accessed")
	w.WriteHeader(http.StatusOK)
}

// login Used to authenticate a user based on email and hashed password. Generated JWT token and
// places it in the cookie content.
func login(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		log.Errorf("login - error in json decoding: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	db, err := internal.New(dbConfig)

	account, err := db.Login(credentials.Email, credentials.HashedPassword)
	if err != nil {
		log.Errorf("login - error when retreving account: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if account == nil {
		log.Errorf("login - account is nil: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	cookieContent, err := generateJWT(credentials)
	if err != nil {
		log.Errorf("login - error while generating cookie content: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "token", Value: cookieContent.Token,
		Expires:  cookieContent.ExpirationTime,
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
	return
}

// registerAccount Used to register a new account based on the request.
func registerAccount(w http.ResponseWriter, r *http.Request) {
	newAccount := types.Account{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("registerAccount - error when reading request body: %v", err)
	}

	err = json.Unmarshal(body, &newAccount)
	if err != nil {
		log.Errorf("registerAccount - error when trying to unmarshal request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	db, err := internal.New(dbConfig)

	if err != nil {
		log.Errorf("registerAccount - error when trying to connect to database: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.CreateAccount(&newAccount)
	if err != nil {
		log.Errorf("registerAccount - error when trying to create a new account: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// updateAccount Used to update an existing account based on the request.
func updateAccount(w http.ResponseWriter, r *http.Request) {
	updatedAccount := types.Account{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("updateAccount - error when reading request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &updatedAccount)
	if err != nil {
		log.Errorf("updateAccount - error when trying to unmarshal request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	db, err := internal.New(dbConfig)

	if err != nil {
		log.Errorf("updateAccount - error when trying to connect to database: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.UpdateAccount(&updatedAccount)
	if err != nil {
		log.Errorf("updateAccount - error when trying to update an account: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// removeAccount Used to remove an account based on the request.
func removeAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		log.Errorf("removeAccount - error when parsing id parameter: not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, err := internal.New(dbConfig)

	if err != nil {
		log.Errorf("removeAccount - error when trying to connect to database: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.DeleteAccount(id)
	if err != nil {
		log.Errorf("removeAccount - error when trying to delete the account: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// generateJWT Used to generate a cookie content, that contains a JWT and expirationTime.
func generateJWT(credentials Credentials) (*CookieContent, error) {
	expirationTime := time.Now().Add(JwtExpirationTime * time.Minute)
	claims := &Claims{
		Username: credentials.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    os.Getenv("ISSUER"),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, fmt.Errorf("generateJWT - error when trying to sign token: %v", err)
	}

	return &CookieContent{
		Token:          tokenString,
		ExpirationTime: expirationTime,
	}, nil
}

// validateJWT Used to validate the JWT present in the cookie.
func validateJWT(request *http.Request) error {
	cookie, err := request.Cookie("token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return fmt.Errorf("validateJWT - error when no cookie is set: %v", err)
		}
		return fmt.Errorf("validateJWT - error when checking cookie content: %v", err)
	}

	token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return fmt.Errorf("validateJWT - error when validatating cookie signature: %v", err)
		}

		return fmt.Errorf("validateJWT - error when parsing cookie: %v", err)
	}

	if !token.Valid {
		return fmt.Errorf("invalid CookieContent")
	}

	return nil
}

// LoadEnvironmentVariable Used to load constants from the environment variables.
func LoadEnvironmentVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
