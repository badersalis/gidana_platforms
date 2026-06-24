// @title           Gidana API
// @version         1.0.0
// @description     Real-estate platform API — property listings, rentals, messaging, digital wallets & payments.
// @contact.name    Gidana Support
// @contact.email   support@gidana.com
// @host            localhost:8080
// @BasePath        /api/v1
// @schemes         http https

// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     Type "Bearer" followed by a space and your JWT token.

package main

import (
	"log"
	"os"

	_ "github.com/badersalis/gidana_backend/docs"
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
	log.Printf("Scalar UI  → http://localhost:%s/docs", port)
	log.Printf("Swagger UI → http://localhost:%s/swagger/index.html", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
