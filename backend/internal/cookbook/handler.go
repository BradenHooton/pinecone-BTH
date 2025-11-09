package cookbook

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/BradenHooton/pinecone-api/internal/middleware"
	"github.com/BradenHooton/pinecone-api/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for cookbooks
type Handler struct {
	service *Service
}

// NewHandler creates a new cookbook handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// HandleCreate handles POST /api/v1/cookbooks
func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.CreateCookbookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cookbook, err := h.service.CreateCookbook(r.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to create cookbook")
		return
	}

	writeJSON(w, http.StatusCreated, models.CookbookResponse{Data: *cookbook})
}

// HandleList handles GET /api/v1/cookbooks
func (h *Handler) HandleList(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse pagination params
	limit := int64(20)
	offset := int64(0)

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 64); err == nil {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.ParseInt(offsetStr, 10, 64); err == nil {
			offset = o
		}
	}

	cookbooks, total, err := h.service.GetCookbooksByUser(r.Context(), userID, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get cookbooks")
		return
	}

	response := models.CookbookListResponse{
		Data: cookbooks,
	}
	response.Meta.Total = total

	writeJSON(w, http.StatusOK, response)
}

// HandleGetByID handles GET /api/v1/cookbooks/{id}
func (h *Handler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get cookbook ID from URL
	idStr := chi.URLParam(r, "id")
	cookbookID, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid cookbook ID")
		return
	}

	cookbook, err := h.service.GetCookbookByID(r.Context(), userID, cookbookID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "Cookbook not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			writeError(w, http.StatusForbidden, "Forbidden")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get cookbook")
		return
	}

	writeJSON(w, http.StatusOK, models.CookbookResponse{Data: *cookbook})
}

// HandleUpdate handles PUT /api/v1/cookbooks/{id}
func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get cookbook ID from URL
	idStr := chi.URLParam(r, "id")
	cookbookID, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid cookbook ID")
		return
	}

	var req models.UpdateCookbookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cookbook, err := h.service.UpdateCookbook(r.Context(), userID, cookbookID, &req)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "Cookbook not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			writeError(w, http.StatusForbidden, "Forbidden")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to update cookbook")
		return
	}

	writeJSON(w, http.StatusOK, models.CookbookResponse{Data: *cookbook})
}

// HandleDelete handles DELETE /api/v1/cookbooks/{id}
func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get cookbook ID from URL
	idStr := chi.URLParam(r, "id")
	cookbookID, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid cookbook ID")
		return
	}

	if err := h.service.DeleteCookbook(r.Context(), userID, cookbookID); err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "Cookbook not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			writeError(w, http.StatusForbidden, "Forbidden")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to delete cookbook")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleAddRecipe handles POST /api/v1/cookbooks/{cookbook_id}/recipes/{recipe_id}
func (h *Handler) HandleAddRecipe(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get cookbook ID and recipe ID from URL
	cookbookIDStr := chi.URLParam(r, "cookbook_id")
	cookbookID, err := uuid.Parse(cookbookIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid cookbook ID")
		return
	}

	recipeIDStr := chi.URLParam(r, "recipe_id")
	recipeID, err := uuid.Parse(recipeIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid recipe ID")
		return
	}

	if err := h.service.AddRecipeToCookbook(r.Context(), userID, cookbookID, recipeID); err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "Cookbook not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			writeError(w, http.StatusForbidden, "Forbidden")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to add recipe to cookbook")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleRemoveRecipe handles DELETE /api/v1/cookbooks/{cookbook_id}/recipes/{recipe_id}
func (h *Handler) HandleRemoveRecipe(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get cookbook ID and recipe ID from URL
	cookbookIDStr := chi.URLParam(r, "cookbook_id")
	cookbookID, err := uuid.Parse(cookbookIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid cookbook ID")
		return
	}

	recipeIDStr := chi.URLParam(r, "recipe_id")
	recipeID, err := uuid.Parse(recipeIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid recipe ID")
		return
	}

	if err := h.service.RemoveRecipeFromCookbook(r.Context(), userID, cookbookID, recipeID); err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "Cookbook not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			writeError(w, http.StatusForbidden, "Forbidden")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to remove recipe from cookbook")
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
