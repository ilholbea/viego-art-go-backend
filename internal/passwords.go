package internal

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
)

func HashPassword(password string) (string, error) {
	cost, err := strconv.Atoi(LoadEnvironmentVariable("HASH_COST"))
	if err != nil {
		return "", err
	}
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(encryptedPassword), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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
