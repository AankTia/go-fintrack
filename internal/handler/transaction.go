package handler

import (
	"encoding/json"
	"errors"
	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/AankTia/go-fintrack/internal/service"
	"net/http"
)

type TransactionHandler struct {
	svc *service.TransactionService
}

func NewTransactionHandler(svc *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{svc: svc}
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	transaction, err := h.svc.CreateTransaction(r.Context(), userID, req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			respondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, transaction)
}
