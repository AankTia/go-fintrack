package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5" // Import chi untuk URLParam
	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/AankTia/go-fintrack/internal/service"
)

type CategoryHandler struct {
	svc *service.CategoryService
}

func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		svc: svc,
	}
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Ambil userID dari context (disuntikkan oleh AuthMiddleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	category, err := h.svc.CreateCategory(r.Context(), userID, req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			respondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, category)
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	categories, err := h.svc.GetMyCategories(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to fetch categories")
		return
	}

	respondWithJSON(w, http.StatusOK, categories)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	categoryID := chi.URLParam(r, "id")
	if categoryID == "" {
		respondWithError(w, http.StatusBadRequest, "category id is required")
		return
	}

	category, err := h.svc.DeleteCategory(r.Context(), categoryID, userID)
	if err != nil {
		if errors.Is(err, service.ErrAccountNotFound) {
			respondWithError(w, http.StatusNotFound, "category not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, category)
}