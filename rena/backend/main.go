package main

import (
	"log"
	"os"
	"strings"

	"rena-platform/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	storage_go "github.com/supabase-community/storage-go"
	"github.com/supabase-community/supabase-go"
)

const defaultPort = "8080"

func loadEnv() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("Warning: Error loading .env file from root")
	}
}

func createSupabaseClients() (*supabase.Client, *storage_go.Client, error) {
	supabaseURL := os.Getenv("VITE_SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
	if supabaseKey == "" {
		supabaseKey = os.Getenv("VITE_SUPABASE_ANON_KEY")
		log.Println("[Debug] WARNING: SUPABASE_SERVICE_ROLE_KEY not found, falling back to VITE_SUPABASE_ANON_KEY")
	}
	if supabaseURL == "" || supabaseKey == "" {
		log.Println("VITE_SUPABASE_URL or SUPABASE key not set in .env")
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
	r.Use(corsMiddleware)
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

func corsMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}
	c.Next()
}
