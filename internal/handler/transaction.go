package handler

import (
	"fmt"
	"encoding/json"
	"errors"
	"log"
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

	if req.Amount <= 0 {
		respondWithError(w, http.StatusBadRequest, "amount must be greater than 0")
		return
	}

	if req.Type != "income" && req.Type != "expense" {
		respondWithError(w, http.StatusBadRequest, "type mus be either 'income' or 'expense'")
		return
	}

	transaction, err := h.svc.CreateTransaction(r.Context(), userID, req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			respondWithError(w, http.StatusBadRequest, err.Error())
		} else if errors.Is(err, service.ErrResourceNotFound) {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			// LOG ERROR ASLI KE CONSOLE
			log.Printf("[ERROR] CreateTransaction failed: %v", err)
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, transaction)
}

func (h *TransactionHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.TransactionListRequest
	req.StartDate = r.URL.Query().Get("start_date")
	req.EndDate = r.URL.Query().Get("end_date")
	req.Type = r.URL.Query().Get("type")
	req.AccountID = r.URL.Query().Get("account_id")

	fmt.Sscanf(r.URL.Query().Get("page"), "%d", &req.Page)
	fmt.Sscanf(r.URL.Query().Get("limit"), "%d", &req.Limit)

	data, meta, err := h.svc.ListTransactions(r.Context(), userID, req)
	if err != nil {
		log.Printf("[ERROR] ListTransactions failed: %v", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data": data,
		"meta": meta,
	})
}

func (h *TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("user_id").(string)
	id := r.URL.Path[len("/api/v1/transactions/"):]

	var req model.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	transaction, err := h.svc.UpdateTransaction(r.Context(), userID, id, req)
	if err != nil {
		if errors.Is(err, service.ErrResourceNotFound) {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			log.Printf("[ERROR] UpdateTransaction failed: %v", err)
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, transaction)
}

func (h *TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("user_id").(string)
	id := r.URL.Path[len("/api/v1/transactions/"):]

	if err := h.svc.DeleteTransaction(r.Context(), userID, id); err != nil {
		if errors.Is(err, service.ErrResourceNotFound) {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			log.Printf("[ERROR] DeleteTransaction failed: %v", err)
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
