package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadDepartments_Success(t *testing.T) {
	// ARRANGE: Create a temporary test YAML file
	yamlContent := `departments:
  - id: produce
    name: Produce
    description: Fresh fruits and vegetables
    order: 1
  - id: meat
    name: Meat & Poultry
    description: Fresh and packaged meats
    order: 2
  - id: dairy
    name: Dairy & Eggs
    description: Milk, cheese, yogurt, eggs
    order: 3
`
	tmpFile, err := os.CreateTemp("", "departments-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(yamlContent)
	require.NoError(t, err)
	tmpFile.Close()

	// ACT
	departments, err := LoadDepartments(tmpFile.Name())

	// ASSERT
	require.NoError(t, err)
	assert.Len(t, departments.Departments, 3)
	assert.Equal(t, "produce", departments.Departments[0].ID)
	assert.Equal(t, "Produce", departments.Departments[0].Name)
	assert.Equal(t, "Fresh fruits and vegetables", departments.Departments[0].Description)
	assert.Equal(t, 1, departments.Departments[0].Order)
}

func TestLoadDepartments_FileNotFound(t *testing.T) {
	// ACT
	_, err := LoadDepartments("nonexistent.yaml")

	// ASSERT
	assert.Error(t, err)
}

func TestLoadDepartments_InvalidYAML(t *testing.T) {
	// ARRANGE
	tmpFile, err := os.CreateTemp("", "invalid-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString("invalid: yaml: content:")
	require.NoError(t, err)
	tmpFile.Close()

	// ACT
	_, err = LoadDepartments(tmpFile.Name())

	// ASSERT
	assert.Error(t, err)
}

func TestGetDepartmentByID_Found(t *testing.T) {
	// ARRANGE
	departments := &DepartmentsConfig{
		Departments: []Department{
			{ID: "produce", Name: "Produce", Order: 1},
			{ID: "meat", Name: "Meat", Order: 2},
		},
	}

	// ACT
	dept, found := departments.GetDepartmentByID("produce")

	// ASSERT
	assert.True(t, found)
	assert.Equal(t, "produce", dept.ID)
	assert.Equal(t, "Produce", dept.Name)
}

func TestGetDepartmentByID_NotFound(t *testing.T) {
	// ARRANGE
	departments := &DepartmentsConfig{
		Departments: []Department{
			{ID: "produce", Name: "Produce", Order: 1},
		},
	}

	// ACT
	_, found := departments.GetDepartmentByID("nonexistent")

	// ASSERT
	assert.False(t, found)
}
