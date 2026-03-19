package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"rena-platform/backend/services"

	"github.com/gin-gonic/gin"
	storage_go "github.com/supabase-community/storage-go"
	"github.com/supabase-community/supabase-go"
)

type AppContext struct {
	Client        *supabase.Client
	StorageClient *storage_go.Client
	DBService     *services.DatabaseService
}

func respondError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

func respondJSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

func getRequiredHeader(c *gin.Context, key string) (string, bool) {
	value := strings.TrimSpace(c.GetHeader(key))
	if value == "" {
		respondError(c, http.StatusUnauthorized, fmt.Sprintf("%s header required", key))
		return "", false
	}
	return value, true
}

func getUserIDFromHeader(c *gin.Context) (string, bool) {
	userID, ok := getRequiredHeader(c, "X-User-ID")
	if !ok {
		return "", false
	}
	if userID == "" {
		respondError(c, http.StatusUnauthorized, "X-User-ID is required")
		return "", false
	}
	return userID, true
}

func (ctx *AppContext) HealthHandler(c *gin.Context) {
	log.Println("[Debug] Health check endpoint called")
	respondJSON(c, http.StatusOK, gin.H{"status": "ok", "message": "Rena Builder Backend is running"})
}

func (ctx *AppContext) ExportProjectHandler(c *gin.Context) {
	projectID := c.Param("id")
	userID, ok := getUserIDFromHeader(c)
	if !ok {
		return
	}

	project, err := ctx.DBService.GetProjectByIDAndUserID(projectID, userID)
	if err != nil {
		log.Printf("[Export] Failed to fetch project: %v\n", err)
		respondError(c, http.StatusNotFound, "Project not found")
		return
	}

	exportData := services.ProjectExport{
		Metadata: services.ProjectMetadata{
			ID:          project.ID,
			Name:        project.Name,
			PackageName: project.PackageName,
			Platform:    project.Platform,
			Color:       project.Color,
			VersionCode: project.VersionCode,
			VersionName: project.VersionName,
			CreatedAt:   project.CreatedAt,
			UpdatedAt:   project.UpdatedAt,
		},
		WorkspaceXML: project.WorkspaceXML,
		Code:         project.GeneratedCode,
	}

	rnpData, err := services.ExportProjectToBytes(exportData)
	if err != nil {
		log.Printf("[Export] Failed to create rnp file: %v\n", err)
		respondError(c, http.StatusInternalServerError, "Failed to create export file")
		return
	}

	filename := fmt.Sprintf("%s.rnp", sanitizeFilename(project.Name))
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Length", fmt.Sprintf("%d", len(rnpData)))
	c.Data(http.StatusOK, "application/zip", rnpData)
	log.Printf("[Export] Successfully exported project: %s\n", filename)
}

func (ctx *AppContext) ImportProjectHandler(c *gin.Context) {
	userID := c.PostForm("user_id")
	if userID == "" {
		respondError(c, http.StatusBadRequest, "user_id is required")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Printf("[Import] Failed to get file: %v\n", err)
		respondError(c, http.StatusBadRequest, "Failed to get uploaded file")
		return
	}
	defer file.Close()

	fileBuffer := make([]byte, header.Size)
	_, err = file.Read(fileBuffer)
	if err != nil {
		log.Printf("[Import] Failed to read file: %v\n", err)
		respondError(c, http.StatusInternalServerError, "Failed to read file")
		return
	}

	exportData, err := services.ImportProjectFromBytes(fileBuffer)
	if err != nil {
		log.Printf("[Import] Failed to parse rnp file: %v\n", err)
		respondError(c, http.StatusBadRequest, "Invalid .rnp file format")
		return
	}

	newProject := map[string]interface{}{
		"user_id":        userID,
		"name":           exportData.Metadata.Name + " (Imported)",
		"package_name":   exportData.Metadata.PackageName,
		"platform":       exportData.Metadata.Platform,
		"color":          exportData.Metadata.Color,
		"version_code":   exportData.Metadata.VersionCode,
		"version_name":   exportData.Metadata.VersionName,
		"workspace_xml":  exportData.WorkspaceXML,
		"generated_code": exportData.Code,
	}

	respBody, _, err := ctx.Client.From("projects").Insert(newProject, false, "", "", "merge-duplicates").Execute()
	if err != nil {
		log.Printf("[Import] Failed to save imported project: %v\n", err)
		respondError(c, http.StatusInternalServerError, "Failed to save imported project")
		return
	}

	log.Printf("[Import] Successfully imported project: %s\n", exportData.Metadata.Name)
	respondJSON(c, http.StatusOK, gin.H{"message": "Project imported successfully", "project_name": exportData.Metadata.Name + " (Imported)", "response": string(respBody)})
}

func (ctx *AppContext) GetUserHandler(c *gin.Context) {
	userID, ok := getRequiredHeader(c, "X-User-ID")
	if !ok {
		return
	}

	user, err := ctx.DBService.GetUserByID(userID)
	if err != nil {
		log.Printf("[Auth] Failed to get user: %v\n", err)
		respondError(c, http.StatusNotFound, "User not found")
		return
	}

	respondJSON(c, http.StatusOK, gin.H{"user": user})
}

func (ctx *AppContext) GetProjectsHandler(c *gin.Context) {
	userID, ok := getRequiredHeader(c, "X-User-ID")
	if !ok {
		return
	}

	projects, err := ctx.DBService.GetProjectsByUserID(userID)
	if err != nil {
		log.Printf("[Projects] Failed to fetch projects: %v\n", err)
		respondError(c, http.StatusInternalServerError, "Failed to fetch projects")
		return
	}

	respondJSON(c, http.StatusOK, gin.H{"projects": projects})
}

func (ctx *AppContext) GetProjectHandler(c *gin.Context) {
	projectID := c.Param("id")
	userID, ok := getUserIDFromHeader(c)
	if !ok {
		return
	}

	project, err := ctx.DBService.GetProjectByIDAndUserID(projectID, userID)
	if err != nil {
		log.Printf("[Projects] Failed to fetch project: %v\n", err)
		respondError(c, http.StatusNotFound, "Project not found")
		return
	}
	respondJSON(c, http.StatusOK, gin.H{"project": project})
}

func (ctx *AppContext) CreateProjectHandler(c *gin.Context) {
	userID, ok := getUserIDFromHeader(c)
	if !ok {
		return
	}

	var req struct {
		Name          string `json:"name"`
		PackageName   string `json:"package_name"`
		Platform      string `json:"platform"`
		Color         string `json:"color"`
		IconURL       string `json:"icon_url"`
		VersionCode   int    `json:"version_code"`
		VersionName   string `json:"version_name"`
		WorkspaceXML  string `json:"workspace_xml"`
		GeneratedCode string `json:"generated_code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if userID == "" || req.Name == "" || req.PackageName == "" {
		respondError(c, http.StatusBadRequest, "Missing required fields")
		return
	}

	project := map[string]interface{}{
		"user_id":        userID,
		"name":           req.Name,
		"package_name":   req.PackageName,
		"platform":       req.Platform,
		"color":          req.Color,
		"icon_url":       req.IconURL,
		"version_code":   req.VersionCode,
		"version_name":   req.VersionName,
		"workspace_xml":  req.WorkspaceXML,
		"generated_code": req.GeneratedCode,
	}

	createdProject, err := ctx.DBService.CreateProject(project)
	if err != nil {
		log.Printf("[Projects] Failed to create project: %v\n", err)
		respondError(c, http.StatusInternalServerError, "Failed to create project")
		return
	}

	respondJSON(c, http.StatusOK, gin.H{"project": createdProject, "message": "Project created successfully"})
}

func (ctx *AppContext) UpdateProjectHandler(c *gin.Context) {
	projectID := c.Param("id")

	var req struct {
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
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.PackageName != "" {
		updates["package_name"] = req.PackageName
	}
	if req.Platform != "" {
		updates["platform"] = req.Platform
	}
	if req.Color != "" {
		updates["color"] = req.Color
	}
	if req.IconURL != "" {
		updates["icon_url"] = req.IconURL
	}
	if req.VersionCode > 0 {
		updates["version_code"] = req.VersionCode
	}
	if req.VersionName != "" {
		updates["version_name"] = req.VersionName
	}
	if req.WorkspaceXML != "" {
		updates["workspace_xml"] = req.WorkspaceXML
	}
	if req.GeneratedCode != "" {
		updates["generated_code"] = req.GeneratedCode
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	updates["updated_at"] = "NOW()"

	userID, ok := getUserIDFromHeader(c)
	if !ok {
		return
	}

	updatedProject, err := ctx.DBService.UpdateProject(projectID, userID, updates)
	if err != nil {
		log.Printf("[Projects] Failed to update project: %v\n", err)
		respondError(c, http.StatusInternalServerError, "Failed to update project")
		return
	}

	respondJSON(c, http.StatusOK, gin.H{"project": updatedProject, "message": "Project updated successfully"})
}

func (ctx *AppContext) DeleteProjectHandler(c *gin.Context) {
	projectID := c.Param("id")
	userID, ok := getUserIDFromHeader(c)
	if !ok {
		return
	}
	if err := ctx.DBService.DeleteProject(projectID, userID); err != nil {
		log.Printf("[Projects] Failed to delete project: %v\n", err)
		respondError(c, http.StatusInternalServerError, "Failed to delete project")
		return
	}
	respondJSON(c, http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

func (ctx *AppContext) GenerateKeystoreHandler(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id"`
		Email  string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("[Keystore] Invalid request body:", err)
		respondError(c, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.Email == "" {
		respondError(c, http.StatusBadRequest, "Email is required")
		return
	}

	parts := strings.Split(req.Email, "@")
	username := parts[0]
	derivedPassword := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(username, " ", ""), ",", ""), ".", "")
	if len(derivedPassword) < 6 {
		derivedPassword += "rena123"
	}

	keystoreName := req.UserID + ".keystore"
	keystorePath, err := services.GenerateKeystore(req.UserID, derivedPassword)
	if err != nil {
		log.Println("[Keystore] Generation failed:", err)
		respondError(c, http.StatusInternalServerError, "Failed to generate keystore: "+err.Error())
		return
	}

	keystoreData, err := os.ReadFile(keystorePath)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "Failed to read keystore data")
		return
	}

	_, err = ctx.StorageClient.UploadFile("keystores", keystoreName, bytes.NewReader(keystoreData))
	if err != nil {
		if strings.Contains(err.Error(), "409") || strings.Contains(err.Error(), "already exists") {
			log.Printf("[Keystore] Keystore %s already exists in storage, skipping upload.\n", keystoreName)
		} else {
			log.Printf("[Keystore] Storage upload WARNING: %v\n", err)
		}
	} else {
		log.Println("[Keystore] Storage upload successful")
	}

	keystoreMetadata := map[string]interface{}{"user_id": req.UserID, "keystore_name": keystoreName, "alias": "rena_key"}
	respBody, count, err := ctx.Client.From("user_keystores").Upsert(keystoreMetadata, "user_id", "", "merge-duplicates").Execute()
	if err != nil {
		log.Printf("[Keystore] Database upsert FAILED: %v, Response: %s, Count: %d\n", err, string(respBody), count)
	} else {
		log.Printf("[Keystore] Database metadata saved successfully. Response: %s\n", string(respBody))
	}

	respondJSON(c, http.StatusOK, gin.H{"message": "Keystore processed successfully", "name": keystoreName})
}

func getString(m map[string]interface{}, key string, defaultVal ...string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return ""
}

func getInt(m map[string]interface{}, key string, defaultVal ...int) int {
	if v, ok := m[key].(float64); ok {
		return int(v)
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return 0
}

func sanitizeFilename(filename string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	return re.ReplaceAllString(filename, "_")
}
