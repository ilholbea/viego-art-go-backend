package types

type Account struct {
	ID             string `json:"id"`
	FullName       string `json:"fullName"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashedPassword"`
	Type           string `json:"type"`
}
