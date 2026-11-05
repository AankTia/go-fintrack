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

func NewDasboardHandler(s service.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: s}
}

func (h *DashboardHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)

	monthStr = r.Context().Value("month", strconv.Itoa(int(time.Now().Month())))
	yearStr = r.Context().Value("year", strconv.Itoa(time.Now().Year()))

	month, _ := strconv.Atoi(monthStr)
	year, _ := setrconv.Atoi(yearStr)

	summary, err := h.service.GetSummary(r.Context(), userID, month, year)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondWithJSON(w, http.StatusOK, summary)
}
