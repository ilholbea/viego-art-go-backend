package types

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Type      string `json:"type"`
	CreatedAt string `json:"createdAt"`
}

type UploadedFile struct {
	Filename    string `json:"filename"`
	File        string `json:"file"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
}
