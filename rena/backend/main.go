package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"rena-platform/backend/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	storage_go "github.com/supabase-community/storage-go"
	"github.com/supabase-community/supabase-go"
)

const defaultPort = "8080"

func loadEnv() {
	// In production (Railway), env vars come from platform secrets, not .env.
	if os.Getenv("RAILWAY_ENVIRONMENT") != "" || os.Getenv("GO_ENV") == "production" {
		log.Println("[Debug] Running in production; skipping .env file load")
		return
	}

	paths := []string{".env", "../.env", "../../.env", "../../../.env"}
	for _, p := range paths {
		err := godotenv.Load(p)
		if err == nil {
			log.Printf("[Debug] Loaded env from: %s\n", p)
			return
		}
	}
	log.Println("[Debug] .env not found; continuing with runtime environment variables")
}

func createSupabaseClients() (*supabase.Client, *storage_go.Client, error) {
	supabaseURL := strings.TrimSpace(os.Getenv("VITE_SUPABASE_URL"))
	supabaseKey := strings.TrimSpace(os.Getenv("SUPABASE_SERVICE_ROLE_KEY"))
	if supabaseKey == "" {
		supabaseKey = strings.TrimSpace(os.Getenv("VITE_SUPABASE_ANON_KEY"))
		log.Println("[Debug] SUPABASE_SERVICE_ROLE_KEY missing: falling back to VITE_SUPABASE_ANON_KEY")
	}

	if supabaseURL == "" || supabaseKey == "" {
		log.Println("ERROR: Missing required Supabase settings. Set VITE_SUPABASE_URL and SUPABASE_SERVICE_ROLE_KEY or VITE_SUPABASE_ANON_KEY.")
		log.Printf("Current env: VITE_SUPABASE_URL=%q, SUPABASE_SERVICE_ROLE_KEY set=%v, VITE_SUPABASE_ANON_KEY set=%v\n",
			supabaseURL,
			supabaseKey != "",
			os.Getenv("VITE_SUPABASE_ANON_KEY") != "",
		)
		return nil, nil, fmt.Errorf("missing supabase config environment variables")
	}

	client, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		return nil, nil, err
	}
	storageClient := storage_go.NewClient(supabaseURL+"/storage/v1", supabaseKey, nil)
	return client, storageClient, nil
}

func initStorageBucket(storageClient *storage_go.Client, bucketName string) {
	_, err := storageClient.GetBucket(bucketName)
	if err != nil {
		if strings.Contains(err.Error(), "Invalid Compact JWS") {
			log.Println("[Debug] WARNING: Invalid Compact JWS. Is SUPABASE_SERVICE_ROLE_KEY correct?")
			return
		}
		log.Printf("[Debug] Bucket '%s' not found; creating...\n", bucketName)
		_, err = storageClient.CreateBucket(bucketName, storage_go.BucketOptions{Public: false})
		if err != nil {
			log.Printf("[Debug] WARNING: Could not create bucket '%s': %v\n", bucketName, err)
			return
		}
	}
}

func main() {
	loadEnv()
	client, storageClient, err := createSupabaseClients()
	if err != nil {
		log.Fatal("Failed to create Supabase client:", err)
	}
	dbService := services.NewDatabaseService(client)
	initStorageBucket(storageClient, "keystores")
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://aemo-dev.github.io", "http://localhost:8080", "https://rena-backend-production-e6a7.up.railway.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-User-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 24 * time.Hour,
	}))
	r.SetTrustedProxies(nil)
	appCtx := &AppContext{Client: client, StorageClient: storageClient, DBService: dbService}
	SetupRoutes(r, appCtx)
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	log.Printf("Rena Builder Backend starting on :%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}

