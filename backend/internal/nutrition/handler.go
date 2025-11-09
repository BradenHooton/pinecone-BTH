package nutrition

import (
	"encoding/json"
	"net/http"

	"github.com/BradenHooton/pinecone-api/internal/models"
)

// Handler handles HTTP requests for nutrition
type Handler struct {
	service *Service
}

// NewHandler creates a new nutrition handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// HandleSearch handles GET /api/v1/nutrition/search?query={query}
func (h *Handler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		writeError(w, http.StatusBadRequest, "Query parameter is required")
		return
	}

	results, err := h.service.Search(r.Context(), query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to search nutrition data")
		return
	}

	response := models.NutritionSearchResponse{
		Data: results,
	}
	response.Meta.Total = len(results)

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
