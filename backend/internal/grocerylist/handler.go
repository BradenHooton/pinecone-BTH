package grocerylist

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

// Handler handles HTTP requests for grocery lists
type Handler struct {
	service *Service
}

// NewHandler creates a new grocery list handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// HandleCreate handles POST /api/v1/grocery-lists
func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.CreateGroceryListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	groceryList, err := h.service.CreateGroceryList(r.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to create grocery list")
		return
	}

	writeJSON(w, http.StatusCreated, models.GroceryListResponse{Data: *groceryList})
}

// HandleList handles GET /api/v1/grocery-lists
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

	groceryLists, total, err := h.service.GetGroceryListsByUser(r.Context(), userID, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get grocery lists")
		return
	}

	response := models.GroceryListListResponse{
		Data: groceryLists,
	}
	response.Meta.Total = total

	writeJSON(w, http.StatusOK, response)
}

// HandleGetByID handles GET /api/v1/grocery-lists/{id}
func (h *Handler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get grocery list ID from URL
	idStr := chi.URLParam(r, "id")
	groceryListID, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid grocery list ID")
		return
	}

	groceryList, err := h.service.GetGroceryListByID(r.Context(), userID, groceryListID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "Grocery list not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			writeError(w, http.StatusForbidden, "Forbidden")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get grocery list")
		return
	}

	writeJSON(w, http.StatusOK, models.GroceryListResponse{Data: *groceryList})
}

// HandleDelete handles DELETE /api/v1/grocery-lists/{id}
func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get grocery list ID from URL
	idStr := chi.URLParam(r, "id")
	groceryListID, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid grocery list ID")
		return
	}

	if err := h.service.DeleteGroceryList(r.Context(), userID, groceryListID); err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "Grocery list not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			writeError(w, http.StatusForbidden, "Forbidden")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to delete grocery list")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleAddManualItem handles POST /api/v1/grocery-lists/{id}/items
func (h *Handler) HandleAddManualItem(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get grocery list ID from URL
	idStr := chi.URLParam(r, "id")
	groceryListID, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid grocery list ID")
		return
	}

	var req models.CreateManualItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	item, err := h.service.AddManualItem(r.Context(), userID, groceryListID, &req)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "Grocery list not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			writeError(w, http.StatusForbidden, "Forbidden")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to add manual item")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{"data": item})
}

// HandleUpdateItemStatus handles PATCH /api/v1/grocery-lists/items/{item_id}
func (h *Handler) HandleUpdateItemStatus(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get item ID from URL
	itemIDStr := chi.URLParam(r, "item_id")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	var req models.UpdateItemStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.UpdateItemStatus(r.Context(), userID, itemID, &req); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update item status")
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
