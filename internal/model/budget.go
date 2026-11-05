package model

import "time"

type Budget struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	CategoryID string    `json:"category_id"`
	Amount     int64     `json:"amount"`
	Period     string    `json:"period"` // Format: YYYY-MM
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateBudgetRequest struct {
	CategoryID string `json:"category_id" validate:"required"`
	Amount     int64  `json:"amount" validate:"required,gt=0"`
	Period     string `json:"period" validate:"required,len=7"` // YYYY-MM
}
