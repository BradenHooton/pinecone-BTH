package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Department represents a grocery store department
type Department struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Order       int    `yaml:"order"`
}

// DepartmentsConfig holds all grocery departments
type DepartmentsConfig struct {
	Departments []Department `yaml:"departments"`
}

// LoadDepartments reads the grocery departments from a YAML file
func LoadDepartments(filePath string) (*DepartmentsConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read departments file: %w", err)
	}

	var config DepartmentsConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse departments YAML: %w", err)
	}

	return &config, nil
}

// GetDepartmentByID finds a department by its ID
func (dc *DepartmentsConfig) GetDepartmentByID(id string) (Department, bool) {
	for _, dept := range dc.Departments {
		if dept.ID == id {
			return dept, true
		}
	}
	return Department{}, false
}

// GetAllDepartmentIDs returns a list of all department IDs
func (dc *DepartmentsConfig) GetAllDepartmentIDs() []string {
	ids := make([]string, 0, len(dc.Departments))
	for _, dept := range dc.Departments {
		ids = append(ids, dept.ID)
	}
	return ids
}
