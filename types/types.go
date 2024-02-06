package types

type Account struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Type     string `json:"type"`
}
