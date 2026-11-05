package model

import "time"

type Transaction struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	AccountId       string    `json:"account_id"`
	CategoryId      string    `json:"category_id"`
	Amount          int64     `json:"amount"`
	Type            string    `json:"type"` // income, expense
	Description     string    `json:"description"`
	TransactionDate time.Time `json:"transaction_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateTransactionRequest struct {
	AccountId   string `json:"account_id" validate:"required"`
	CategoryId  string `json:"category_id" validate:"required"`
	Amount      int64  `json:"amount" validate:"required,gt=0"`
	Type        string `json:"type" validate:"required,oneof=income expense"`
	Description string `json:"description"`
	Date        string `json:"date"` // Format YYY-MM-DD
}
