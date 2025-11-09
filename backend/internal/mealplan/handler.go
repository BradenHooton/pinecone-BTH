package mealplan

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/BradenHooton/pinecone-api/internal/models"
)

// Handler handles HTTP requests for meal plans
type Handler struct {
	service *Service
}

// NewHandler creates a new meal plan handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// HandleGetByDate handles GET /api/v1/meal-plans/{date}
func (h *Handler) HandleGetByDate(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		writeError(w, http.StatusBadRequest, "Date parameter is required (format: YYYY-MM-DD)")
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD")
		return
	}

	mealPlan, err := h.service.GetMealPlanByDate(r.Context(), date)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "Meal plan not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get meal plan")
		return
	}

	writeJSON(w, http.StatusOK, models.MealPlanResponse{Data: *mealPlan})
}

// HandleGetByDateRange handles GET /api/v1/meal-plans?start_date={date}&end_date={date}
func (h *Handler) HandleGetByDateRange(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if startDateStr == "" || endDateStr == "" {
		writeError(w, http.StatusBadRequest, "start_date and end_date parameters are required (format: YYYY-MM-DD)")
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid start_date format. Use YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid end_date format. Use YYYY-MM-DD")
		return
	}

	mealPlans, err := h.service.GetMealPlansByDateRange(r.Context(), startDate, endDate)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get meal plans")
		return
	}

	response := models.MealPlanListResponse{
		Data: mealPlans,
	}
	response.Meta.StartDate = startDateStr
	response.Meta.EndDate = endDateStr

	writeJSON(w, http.StatusOK, response)
}

// HandleUpdate handles PUT /api/v1/meal-plans/{date}
func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		writeError(w, http.StatusBadRequest, "Date parameter is required (format: YYYY-MM-DD)")
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD")
		return
	}

	var req models.UpdateMealPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	mealPlan, err := h.service.UpdateMealPlan(r.Context(), date, &req)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to update meal plan")
		return
	}

	writeJSON(w, http.StatusOK, models.MealPlanResponse{Data: *mealPlan})
}

// Helper functions

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]interface{}{
		"error": map[string]string{
			"message": message,
		},
	})
}
