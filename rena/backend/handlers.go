package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"rena-platform/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AppContext struct {
	DB *services.DatabaseService
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type authResponse struct {
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

type deviceAuthResponse struct {
	UserID     int64  `json:"user_id"`
	Token      string `json:"token"`
	Email      string `json:"email"`
	DeviceName string `json:"device_name"`
}

func (app *AppContext) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing, err := app.DB.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	user, err := app.DB.CreateUser(req.Email, string(hashed))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}

	token, err := createJWTToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
		return
	}
	c.JSON(http.StatusCreated, authResponse{UserID: user.ID, Token: token})
}

func (app *AppContext) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := app.DB.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := createJWTToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
		return
	}
	c.JSON(http.StatusOK, authResponse{UserID: user.ID, Token: token})
}

func sanitizeDeviceName(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return "device"
	}

	var b strings.Builder
	lastWasDash := false
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			lastWasDash = false
		case r == '-' || r == '_' || r == '.' || r == ' ':
			if !lastWasDash {
				b.WriteByte('-')
				lastWasDash = true
			}
		}
	}

	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "device"
	}
	return out
}

func (app *AppContext) DeviceLogin(c *gin.Context) {
	hostname, err := os.Hostname()
	if err != nil || strings.TrimSpace(hostname) == "" {
		hostname = "device"
	}

	deviceName := strings.TrimSpace(hostname)
	safeName := sanitizeDeviceName(deviceName)
	email := fmt.Sprintf("%s@device.local", safeName)

	user, err := app.DB.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	if user == nil {
		user, err = app.DB.CreateUser(email, "device-local-user")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create device user"})
			return
		}
	}

	token, err := createJWTToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
		return
	}

	c.JSON(http.StatusOK, deviceAuthResponse{
		UserID:     user.ID,
		Token:      token,
		Email:      email,
		DeviceName: deviceName,
	})
}

func (app *AppContext) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		parsed, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte("supersecretkey"), nil
		})
		if err != nil || !parsed.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := parsed.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		sub, ok := claims["sub"]
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		var userID int64
		switch v := sub.(type) {
		case float64:
			userID = int64(v)
		case int64:
			userID = v
		case string:
			userID, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
				return
			}
		default:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

func (app *AppContext) CreateProject(c *gin.Context) {
	user, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}
	userID := user.(int64)

	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p := &services.Project{
		UserID:       userID,
		Name:         req.Title,
		PackageName:  "com.rena.project",
		Platform:     "android",
		WorkspaceXML: req.Content,
		Status:       "draft",
	}

	out, err := app.DB.CreateProject(p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create project"})
		return
	}
	c.JSON(http.StatusCreated, out)
}

func (app *AppContext) ListProjects(c *gin.Context) {
	user, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}
	userID := user.(int64)

	projects, err := app.DB.GetProjectsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not list projects"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

func (app *AppContext) GetProject(c *gin.Context) {
	user, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}
	userID := user.(int64)

	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	project, err := app.DB.GetProjectByIDAndUserID(projectID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}
	c.JSON(http.StatusOK, project)
}

func (app *AppContext) UpdateProject(c *gin.Context) {
	user, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}
	userID := user.(int64)

	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := app.DB.UpdateProject(projectID, userID, map[string]interface{}{
		"name":          req.Title,
		"workspace_xml": req.Content,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update project"})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (app *AppContext) DeleteProject(c *gin.Context) {
	user, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}
	userID := user.(int64)

	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	if err := app.DB.DeleteProject(projectID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete project"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (app *AppContext) ExportProject(c *gin.Context) {
	user, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}
	userID := user.(int64)

	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	project, err := app.DB.GetProjectByIDAndUserID(projectID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	data, _ := json.Marshal(project)
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	f, _ := zw.Create(filepath.Base(fmt.Sprintf("project-%d.json", project.ID)))
	_, _ = f.Write(data)
	_ = zw.Close()

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=project-%d.zip", project.ID))
	c.Data(http.StatusOK, "application/zip", buf.Bytes())
}

func (app *AppContext) ImportProject(c *gin.Context) {
	user, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}
	userID := user.(int64)

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing file"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not open file"})
		return
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read file"})
		return
	}

	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid zip file"})
		return
	}

	if len(zr.File) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "zip has no files"})
		return
	}

	zf := zr.File[0]
	rc, err := zf.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not open inner file"})
		return
	}
	defer rc.Close()

	projectBytes, err := io.ReadAll(rc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not read project"})
		return
	}

	var project services.Project
	if err := json.Unmarshal(projectBytes, &project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project format"})
		return
	}

	project.UserID = userID
	project.ID = 0
	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()
	created, err := app.DB.CreateProject(&project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not import project"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"project_id": created.ID})
}

func createJWTToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("supersecretkey"))
}
