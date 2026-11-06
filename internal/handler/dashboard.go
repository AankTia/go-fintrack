package handler

import (
	"fmt"
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
	fmt.Println("DEBUG: Entering GetSummary")

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		fmt.Println("DEBUG: UserID not found in context")
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	
	monthVal := r.Context().Value("month")
	fmt.Printf("DEBUG: monthVal type: %T, value: %v\n", monthVal, monthVal)

	var monthStr string
	if monthVal == nil {
		monthStr = strconv.Itoa(int(time.Now().Month()))
	} else {
		monthStr = monthVal.(string)
	}

	yearVal := r.Context().Value("year")
	fmt.Printf("DEBUG: yearVal type: %T, value: %v\n", yearVal, yearVal)

	var yearStr string
	if yearVal == nil {
		yearStr = strconv.Itoa(time.Now().Year())
	} else {
		fmt.Printf("DEBUG: yearVal is NOT a string, it is %T\n", yearVal)

		yearStr = yearVal.(string)
	}

	month, _ := strconv.Atoi(monthStr)
	year, _ := strconv.Atoi(yearStr)

	fmt.Printf("DEBUG: Calling service.GetSummary with userID: %s, month: %d, year: %d\n", userID, month, year)
	summary, err := h.service.GetSummary(r.Context(), userID, month, year)
	if err != nil {
		fmt.Printf("DEBUG: Service error: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	fmt.Println("DEBUG: Success! Sending response")
	respondWithJSON(w, http.StatusOK, summary)
}

func (h *DashboardHandler) GetBreakdown(w http.ResponseWriter, r *http.Request) {
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

	breakdown, err := h.service.GetBreakdown(r.Context(), userID, month, year)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondWithJSON(w, http.StatusOK, breakdown)
}

func (h *DashboardHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value("user_id").(string)
    if !ok {
        respondWithError(w, http.StatusUnauthorized, "unauthorized")
        return
    }

    history, err := h.service.GetHistory(r.Context(), userID)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "internal server error")
        return
    }

    respondWithJSON(w, http.StatusOK, history)
}
