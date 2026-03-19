package main

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine, ctx *AppContext) {
	r.GET("/api/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	r.POST("/api/auth/register", ctx.Register)
	r.POST("/api/auth/login", ctx.Login)

	authGroup := r.Group("/api")
	authGroup.Use(ctx.AuthMiddleware())
	{
		authGroup.GET("/projects", ctx.ListProjects)
		authGroup.GET("/projects/:id", ctx.GetProject)
		authGroup.POST("/projects", ctx.CreateProject)
		authGroup.PUT("/projects/:id", ctx.UpdateProject)
		authGroup.DELETE("/projects/:id", ctx.DeleteProject)
		authGroup.GET("/projects/:id/export", ctx.ExportProject)
		authGroup.POST("/projects/import", ctx.ImportProject)
	}
}
