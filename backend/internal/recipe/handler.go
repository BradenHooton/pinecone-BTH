package recipe

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BradenHooton/pinecone-api/internal/middleware"
	"github.com/BradenHooton/pinecone-api/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const maxUploadSize = 5 * 1024 * 1024 // 5MB

// Handler handles HTTP requests for recipes
type Handler struct {
	service *Service
	uploadDir string
}

// NewHandler creates a new recipe handler
func NewHandler(service *Service, uploadDir string) *Handler {
	return &Handler{
		service: service,
		uploadDir: uploadDir,
	}
}

// HandleCreate handles POST /api/v1/recipes
func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.CreateRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	recipe, err := h.service.CreateRecipe(r.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to create recipe")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"data": recipe,
		"meta": map[string]string{
			"timestamp": recipe.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	})
}

// HandleGetByID handles GET /api/v1/recipes/{id}
func (h *Handler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid recipe ID")
		return
	}

	recipe, err := h.service.GetRecipeByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrRecipeNotFound) {
			writeError(w, http.StatusNotFound, "Recipe not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to get recipe")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data": recipe,
	})
}

// HandleList handles GET /api/v1/recipes
func (h *Handler) HandleList(w http.ResponseWriter, r *http.Request) {
	params := &models.RecipeSearchParams{
		Query:  r.URL.Query().Get("search"),
		Sort:   r.URL.Query().Get("sort"),
		Limit:  20,
		Offset: 0,
	}

	// Parse tags
	if tagsStr := r.URL.Query().Get("tags"); tagsStr != "" {
		params.Tags = strings.Split(tagsStr, ",")
	}

	// Parse limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			params.Limit = limit
		}
	}

	// Parse offset
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			params.Offset = offset
		}
	}

	response, err := h.service.ListRecipes(r.Context(), params)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list recipes")
		return
	}

	writeJSON(w, http.StatusOK, response)
}

// HandleUpdate handles PUT /api/v1/recipes/{id}
func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid recipe ID")
		return
	}

	var req models.UpdateRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	recipe, err := h.service.UpdateRecipe(r.Context(), userID, id, &req)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, ErrRecipeNotFound) {
			writeError(w, http.StatusNotFound, "Recipe not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			writeError(w, http.StatusForbidden, "You don't have permission to update this recipe")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to update recipe")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data": recipe,
	})
}

// HandleDelete handles DELETE /api/v1/recipes/{id}
func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid recipe ID")
		return
	}

	err = h.service.DeleteRecipe(r.Context(), userID, id)
	if err != nil {
		if errors.Is(err, ErrRecipeNotFound) {
			writeError(w, http.StatusNotFound, "Recipe not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			writeError(w, http.StatusForbidden, "You don't have permission to delete this recipe")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to delete recipe")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleUploadImage handles POST /api/v1/recipes/upload-image
func (h *Handler) HandleUploadImage(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form with max size
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		writeError(w, http.StatusBadRequest, "File too large or invalid form data")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Image file is required")
		return
	}
	defer file.Close()

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	allowedTypes := []string{"image/jpeg", "image/jpg", "image/png", "image/webp"}
	if !contains(allowedTypes, contentType) {
		writeError(w, http.StatusBadRequest, "Invalid file type. Allowed: jpg, jpeg, png, webp")
		return
	}

	// Validate file size
	if header.Size > maxUploadSize {
		writeError(w, http.StatusBadRequest, "File size exceeds 5MB limit")
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(h.uploadDir, filename)

	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll(h.uploadDir, 0755); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create upload directory")
		return
	}

	// Save file
	dst, err := os.Create(filePath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to save file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to save file")
		return
	}

	// Return image URL
	imageURL := fmt.Sprintf("/uploads/%s", filename)
	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"data": map[string]string{
			"image_url": imageURL,
		},
	})
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

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
