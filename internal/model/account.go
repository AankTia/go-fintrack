package model

import "time"

// Account adalah entity internal yang merepresentasikan tabel accounts
type Account struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"` // bank, cash, e-wallet
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AccountRequest adalah DTO untuk validasi input saat Create/Update
type AccountRequest struct {
	Name     string `json:"name" validate:"required"`
	Type     string `json:"type" validate:"required,oneof=bank cash e-wallet"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}
