package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/AankTia/go-fintrack/internal/service"
)

type BudgetHandler struct {
	svc service.BudgetService
}

func NewBudgetHandler(svc service.BudgetService) *BudgetHandler {
	return &BudgetHandler{svc: svc}
}

func (h *BudgetHandler) CreateOrUpdate(w http.ResponseWriter, r *http.Request) {
	// Ambil userID dari context (disuntikkan oleh AuthMiddleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.CreateBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	budget, err := h.svc.SetBudget(r.Context(), userID, req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			respondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, budget)
}

func (h *BudgetHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	period := r.URL.Query().Get("period") 
	if period == "" {
		respondWithError(w, http.StatusBadRequest, "period query parameter is required (YYYY-MM)")
		return
	}

	budgets, err := h.svc.GetBudgets(r.Context(), userID, period)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to fetch accounts")
		return
	}

	respondWithJSON(w, http.StatusOK, budgets)
}

func (h *BudgetHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	period := r.URL.Query().Get("period") 
	if period == "" {
		respondWithError(w, http.StatusBadRequest, "period query parameter is required (YYYY-MM)")
		return
	}

	statuses, err := h.svc.GetBudgetStatus(r.Context(), userID, period)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to calculate budget")
		return
	}

	respondWithJSON(w, http.StatusOK, statuses)
}
