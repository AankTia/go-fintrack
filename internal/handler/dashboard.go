package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/AankTia/go-fintrack/internal/service"
)

type DashboardHandler struct {
	service service.DashboardService
}

func NewDashboardHandler(s service.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: s}
}

func (h *DashboardHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	
	monthVal := r.Context().Value("month")
	var monthStr string
	if monthVal == nil {
		monthStr = strconv.Itoa(int(time.Now().Month()))
	} else {
		monthStr = monthVal.(string)
	}

	yearVal := r.Context().Value("year")
	var yearStr string
	if yearVal == nil {
		yearStr = strconv.Itoa(time.Now().Year())
	} else {
		yearStr = yearVal.(string)
	}

	month, _ := strconv.Atoi(monthStr)
	year, _ := strconv.Atoi(yearStr)

	summary, err := h.service.GetSummary(r.Context(), userID, month, year)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondWithJSON(w, http.StatusOK, summary)
}
