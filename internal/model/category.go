package model

import "time"

// Category adalah entity internal yang merepresentasikan tabel categories
type Category struct {
	ID        string    `json:"id"`
	UserID    *string   `json:"user_id"` // Menggunakan pointer karena boleh null (global category)
	Name      string    `json:"name"`
	Type      string    `json:"type"` // income, expense
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CategoryRequest adalah DTO untuk validasi input
type CategoryRequest struct {
	Name string `json:"name" validate:"required"`
	Type string `json:"type" validate:"required,oneof=income expense"`
	Icon string `json:"icon"`
}
