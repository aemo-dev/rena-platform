package main

import (
	"log"
	"os"
	"time"

	"rena-platform/backend/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const defaultPort = "8080"

func loadEnv() {
	if os.Getenv("GO_ENV") == "production" {
		log.Println("[Debug] production mode; skipping .env file load")
		return
	}

	paths := []string{".env", "../.env", "../../.env", "../../../.env"}
	for _, p := range paths {
		err := godotenv.Load(p)
		if err == nil {
			log.Printf("[Debug] loaded env from %s\n", p)
			return
		}
	}
	log.Println("[Debug] no .env found; using runtime env vars")
}

func main() {
	loadEnv()
	db, err := services.NewDatabaseService(os.Getenv("SQLITE_PATH"))
	if err != nil {
		log.Fatal("failed to initialize database:", err)
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:8080", "https://aemo-dev.github.io"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-User-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.SetTrustedProxies(nil)

	appCtx := &AppContext{DB: db}
	SetupRoutes(r, appCtx)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	log.Printf("starting backend on :%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
