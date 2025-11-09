package menu

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/BradenHooton/pinecone-api/internal/models"
)

// Handler handles HTTP requests for menu recommendations
type Handler struct {
	service *Service
}

// NewHandler creates a new menu recommendation handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// HandleRecommend handles POST /api/v1/menu/recommend
func (h *Handler) HandleRecommend(w http.ResponseWriter, r *http.Request) {
	var req models.RecommendRecipesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	recommendations, err := h.service.RecommendRecipes(r.Context(), &req)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get recommendations")
		return
	}

	response := models.RecommendRecipesResponse{
		Data: recommendations,
	}
	response.Meta.ProvidedIngredients = req.Ingredients
	response.Meta.TotalRecipesFound = len(recommendations)

	writeJSON(w, http.StatusOK, response)
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
