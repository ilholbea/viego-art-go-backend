package utils

import (
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	var cost int
	hashCost, ok := os.LookupEnv("HASH_COST")
	if !ok {
		cost = bcrypt.DefaultCost
	}

	cost, err := strconv.Atoi(hashCost)
	if err != nil {
		cost = bcrypt.DefaultCost
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

// ComparePasswords compares a hashed password with its plain-text version
func ComparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil // Return true if they match, false otherwise
}
