package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/AankTia/go-fintrack/internal/service"
)

// RegisterHandler men‑handle POST /api/v1/auth/register
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Panggil service layer
	usr, err := service.RegisterUser(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(usr)
}

// LoginHandler men‑handle POST /api/v1/auth/login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}

	// Panggil service layer
	token, err := service.LoginUser(r.Context(), &req)
	if err != nil {
		// Jika error adalah "invalid credentials", kirim 401 Unauthorized
		if err.Error() == "invalid credentials" {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}

		// Error lainnya kirim 400 atau 500
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
