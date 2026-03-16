package main

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine, ctx *AppContext) {
	// Health check endpoint
	r.GET("/api/health", ctx.HealthHandler)

	// Export/Import
	r.GET("/api/projects/:id/export", ctx.ExportProjectHandler)
	r.POST("/api/projects/import", ctx.ImportProjectHandler)

	// Auth
	r.GET("/api/auth/user", ctx.GetUserHandler)

	// Projects CRUD
	r.GET("/api/projects", ctx.GetProjectsHandler)
	r.GET("/api/projects/:id", ctx.GetProjectHandler)
	r.POST("/api/projects", ctx.CreateProjectHandler)
	r.PUT("/api/projects/:id", ctx.UpdateProjectHandler)
	r.DELETE("/api/projects/:id", ctx.DeleteProjectHandler)

	// Keystore
	r.POST("/api/keystore/generate", ctx.GenerateKeystoreHandler)
}
