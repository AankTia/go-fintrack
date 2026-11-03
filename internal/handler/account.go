package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5" // Import chi untuk URLParam
	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/AankTia/go-fintrack/internal/service"
)

type AccountHandler struct {
	svc *service.AccountService
}

// NewAccountHandler adalah fungsi yang hilang tadi. 
// Harus dimulai dengan huruf Besar agar bisa dipanggil di main.go
func NewAccountHandler(svc *service.AccountService) *AccountHandler {
	return &AccountHandler{
		svc: svc,
	}
}

// Helper untuk mengirim response JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// Helper untuk mengirim error JSON
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Ambil userID dari context (disuntikkan oleh AuthMiddleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.AccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	account, err := h.svc.CreateAccount(r.Context(), userID, req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			respondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, account)
}

func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	accounts, err := h.svc.GetMyAccounts(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to fetch accounts")
		return
	}

	respondWithJSON(w, http.StatusOK, accounts)
}

func (h *AccountHandler) GetDetail(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Mengambil parameter {id} dari URL /api/v1/accounts/{id}
	accountID := chi.URLParam(r, "id")
	if accountID == "" {
		respondWithError(w, http.StatusBadRequest, "account id is required")
		return
	}

	account, err := h.svc.GetAccountDetail(r.Context(), accountID, userID)
	if err != nil {
		if errors.Is(err, service.ErrAccountNotFound) {
			respondWithError(w, http.StatusNotFound, "account not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, account)
}