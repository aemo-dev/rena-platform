package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/supabase-community/supabase-go"
)

// DatabaseService provides database operations through Supabase client
type DatabaseService struct {
	client *supabase.Client
}

// NewDatabaseService creates a new database service instance
func NewDatabaseService(client *supabase.Client) *DatabaseService {
	return &DatabaseService{
		client: client,
	}
}

// User represents a user in the database
type User struct {
	ID        string                 `json:"id"`
	Email     string                 `json:"email"`
	UserMeta  map[string]interface{} `json:"user_metadata"`
	CreatedAt string                 `json:"created_at"`
}

// Project represents a project in the database
type Project struct {
	ID            string `json:"id"`
	UserID        string `json:"user_id"`
	Name          string `json:"name"`
	PackageName   string `json:"package_name"`
	Platform      string `json:"platform"`
	Color         string `json:"color"`
	IconURL       string `json:"icon_url"`
	VersionCode   int    `json:"version_code"`
	VersionName   string `json:"version_name"`
	WorkspaceXML  string `json:"workspace_xml"`
	GeneratedCode string `json:"generated_code"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// GetUserByID retrieves a user by their ID
func (s *DatabaseService) GetUserByID(userID string) (*User, error) {
	resp, count, err := s.client.From("users").Select("*", "", false).Eq("id", userID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if count == 0 {
		return nil, fmt.Errorf("user not found")
	}

	var users []User
	err = parseJSONResponse(resp, &users)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user response: %w", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return &users[0], nil
}

// GetProjectsByUserID retrieves all projects for a specific user
func (s *DatabaseService) GetProjectsByUserID(userID string) ([]Project, error) {
	resp, _, err := s.client.From("projects").Select("*", "", false).Eq("user_id", userID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch projects: %w", err)
	}

	var projects []Project
	err = parseJSONResponse(resp, &projects)
	if err != nil {
		return nil, fmt.Errorf("failed to parse projects response: %w", err)
	}

	return projects, nil
}

// GetProjectByID retrieves a specific project by ID
func (s *DatabaseService) GetProjectByID(projectID string) (*Project, error) {
	resp, count, err := s.client.From("projects").Select("*", "", false).Eq("id", projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch project: %w", err)
	}

	if count == 0 {
		return nil, fmt.Errorf("project not found")
	}

	var projects []Project
	err = parseJSONResponse(resp, &projects)
	if err != nil {
		return nil, fmt.Errorf("failed to parse project response: %w", err)
	}

	if len(projects) == 0 {
		return nil, fmt.Errorf("project not found")
	}

	return &projects[0], nil
}

// CreateProject creates a new project in the database
func (s *DatabaseService) CreateProject(project map[string]interface{}) (*Project, error) {
	resp, _, err := s.client.From("projects").Insert(project, false, "", "", "").Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	var createdProjects []Project
	err = parseJSONResponse(resp, &createdProjects)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created project response: %w", err)
	}

	if len(createdProjects) == 0 {
		return nil, fmt.Errorf("project creation failed")
	}

	return &createdProjects[0], nil
}

// UpdateProject updates an existing project
func (s *DatabaseService) UpdateProject(projectID string, updates map[string]interface{}) (*Project, error) {
	resp, _, err := s.client.From("projects").Update(updates, "id", "").Eq("id", projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	var updatedProjects []Project
	err = parseJSONResponse(resp, &updatedProjects)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updated project response: %w", err)
	}

	if len(updatedProjects) == 0 {
		return nil, fmt.Errorf("project update failed")
	}

	return &updatedProjects[0], nil
}

// DeleteProject deletes a project by ID
func (s *DatabaseService) DeleteProject(projectID string) error {
	_, _, err := s.client.From("projects").Delete("", "").Eq("id", projectID).Execute()
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	return nil
}

// parseJSONResponse parses the Supabase response into the target type
func parseJSONResponse(resp []byte, target interface{}) error {
	// Simple JSON parsing - in production you'd want better error handling
	return json.Unmarshal(resp, target)
}

// Helper function to check if a string slice contains a value
func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// Helper function to sanitize strings
func sanitizeString(s string) string {
	return strings.TrimSpace(s)
}
