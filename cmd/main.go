package main

import (
	"log"
	"os"

	"github.com/badersalis/gidana_backend/internal/config"
	"github.com/badersalis/gidana_backend/internal/database"
	"github.com/badersalis/gidana_backend/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()

	if config.App.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	database.Connect()

	if config.App.SupabaseURL != "" {
		log.Println("Supabase Storage configured")
	} else {
		log.Println("Warning: SUPABASE_URL not set, falling back to local file storage")
	}

	if err := os.MkdirAll(config.App.UploadDir, 0755); err != nil {
		log.Printf("Warning: could not create upload dir: %v", err)
	}

	r := gin.Default()
	r.SetTrustedProxies(nil)

	routes.Setup(r, database.DB)

	port := config.App.Port
	log.Printf("Gidana API server starting on port %s (env: %s)", port, config.App.AppEnv)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
