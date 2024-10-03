package utils

import (
	"context"
	"errors"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/ilholbea/viego-art-backend/types"
	"net/http"
	"os"
	"strings"
	"time"
)

// CustomClaims Custom claims for the JWT that include the user fields
type CustomClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Type     string `json:"type"`
	jwt.RegisteredClaims
}

func CreateToken(user *types.User) (string, error) {
	claims := CustomClaims{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Type:     user.Type,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)), // Token valid for 1 hour
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    os.Getenv("ISSUER"),
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey, ok := os.LookupEnv("JWT_SECRET_KEY")
	if !ok {
		return "", errors.New("JWT_SECRET_KEY environment variable is not set")
	}

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Malformed token", http.StatusUnauthorized)
			return
		}

		// Verify the token and extract claims
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			secretKey, ok := os.LookupEnv("JWT_SECRET_KEY")
			if !ok {
				return nil, errors.New("JWT_SECRET_KEY environment variable is not set")
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Extract the username from the token's claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add the username to the request context
		username := claims["username"].(string)
		ctx := context.WithValue(r.Context(), "username", username)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
