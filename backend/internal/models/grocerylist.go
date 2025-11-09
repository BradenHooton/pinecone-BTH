package models

import (
	"time"

	"github.com/google/uuid"
)

// GroceryDepartment represents grocery store departments
type GroceryDepartment string

const (
	DepartmentProduce   GroceryDepartment = "produce"
	DepartmentMeat      GroceryDepartment = "meat"
	DepartmentSeafood   GroceryDepartment = "seafood"
	DepartmentDairy     GroceryDepartment = "dairy"
	DepartmentBakery    GroceryDepartment = "bakery"
	DepartmentFrozen    GroceryDepartment = "frozen"
	DepartmentPantry    GroceryDepartment = "pantry"
	DepartmentSpices    GroceryDepartment = "spices"
	DepartmentBeverages GroceryDepartment = "beverages"
	DepartmentOther     GroceryDepartment = "other"
)

// GroceryItemStatus represents the status of a grocery item
type GroceryItemStatus string

const (
	StatusPending     GroceryItemStatus = "pending"
	StatusBought      GroceryItemStatus = "bought"
	StatusHaveOnHand  GroceryItemStatus = "have_on_hand"
)

// GroceryList represents a grocery list for a date range
type GroceryList struct {
	ID              uuid.UUID       `json:"id"`
	CreatedByUserID uuid.UUID       `json:"created_by_user_id"`
	StartDate       time.Time       `json:"start_date"`
	EndDate         time.Time       `json:"end_date"`
	Items           []GroceryListItem `json:"items,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// GroceryListItem represents a single item on a grocery list
type GroceryListItem struct {
	ID              uuid.UUID         `json:"id"`
	GroceryListID   uuid.UUID         `json:"grocery_list_id"`
	ItemName        string            `json:"item_name"`
	Quantity        *float64          `json:"quantity,omitempty"`
	Unit            *string           `json:"unit,omitempty"`
	Department      GroceryDepartment `json:"department"`
	Status          GroceryItemStatus `json:"status"`
	IsManual        bool              `json:"is_manual"`
	SourceRecipeID  *uuid.UUID        `json:"source_recipe_id,omitempty"`
	SourceRecipe    *Recipe           `json:"source_recipe,omitempty"`
}

// CreateGroceryListRequest represents the request to create a grocery list
type CreateGroceryListRequest struct {
	StartDate string `json:"start_date"` // YYYY-MM-DD format
	EndDate   string `json:"end_date"`   // YYYY-MM-DD format
}

// CreateManualItemRequest represents the request to add a manual item
type CreateManualItemRequest struct {
	ItemName   string             `json:"item_name"`
	Quantity   *float64           `json:"quantity,omitempty"`
	Unit       *string            `json:"unit,omitempty"`
	Department *GroceryDepartment `json:"department,omitempty"`
}

// UpdateItemStatusRequest represents the request to update an item's status
type UpdateItemStatusRequest struct {
	Status GroceryItemStatus `json:"status"`
}

// GroceryListResponse represents the API response for a grocery list
type GroceryListResponse struct {
	Data GroceryList `json:"data"`
}

// GroceryListListResponse represents the API response for a list of grocery lists
type GroceryListListResponse struct {
	Data []GroceryList `json:"data"`
	Meta struct {
		Total int64 `json:"total"`
	} `json:"meta"`
}

// GroceryItemsByDepartment groups items by department for display
type GroceryItemsByDepartment map[GroceryDepartment][]GroceryListItem
