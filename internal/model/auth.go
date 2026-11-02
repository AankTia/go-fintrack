package model

// RegisterRequest datang dari body POST /auth/register
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest datang dari body POST /auth/login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserResponse men‑return data user setelah register
type UserResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}
