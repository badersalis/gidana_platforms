package routes

import (
	"net/http"
	"strings"

	"github.com/badersalis/gidana_backend/internal/config"
	"github.com/badersalis/gidana_backend/internal/handlers"
	"github.com/badersalis/gidana_backend/internal/middleware"
	"github.com/badersalis/gidana_backend/internal/repositories"
	"github.com/badersalis/gidana_backend/internal/services"
	"github.com/badersalis/gidana_backend/internal/storage"
	appws "github.com/badersalis/gidana_backend/internal/ws"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(r *gin.Engine, db *gorm.DB) {
	allowedOrigins := make(map[string]bool)
	for _, o := range strings.Split(config.App.AllowedOrigins, ",") {
		allowedOrigins[strings.TrimSpace(o)] = true
	}

	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if origin == "" || origin == "null" {
				return true
			}
			return allowedOrigins[origin] || allowedOrigins["*"]
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	r.Static("/uploads", "./uploads")

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "version": "1.0.0"})
	})

	// ── Storage ───────────────────────────────────────────────────────────
	fileStore := storage.NewFileStorage()

	// ── Repositories ──────────────────────────────────────────────────────
	userRepo   := repositories.NewUserRepository(db)
	propRepo   := repositories.NewPropertyRepository(db)
	imageRepo  := repositories.NewPropertyImageRepository(db)
	rentalRepo := repositories.NewRentalRepository(db)
	walletRepo := repositories.NewWalletRepository(db)
	txRepo     := repositories.NewTransactionRepository(db)
	reviewRepo := repositories.NewReviewRepository(db)
	convRepo   := repositories.NewConversationRepository(db)
	msgRepo    := repositories.NewMessageRepository(db)
	favRepo    := repositories.NewFavoriteRepository(db)
	alertRepo  := repositories.NewAlertRepository(db)
	searchRepo := repositories.NewSearchRepository(db)

	// ── Services ──────────────────────────────────────────────────────────
	authSvc   := services.NewAuthService(userRepo)
	userSvc   := services.NewUserService(userRepo, fileStore)
	propSvc   := services.NewPropertyService(propRepo, imageRepo, favRepo, alertRepo, userRepo, fileStore, appws.H)
	rentalSvc := services.NewRentalService(rentalRepo, propRepo)
	walletSvc := services.NewWalletService(walletRepo)
	txSvc     := services.NewTransactionService(walletRepo, txRepo, db)
	reviewSvc := services.NewReviewService(reviewRepo, propRepo)
	msgSvc    := services.NewMessageService(convRepo, msgRepo, userRepo, appws.H)
	favSvc    := services.NewFavoriteService(favRepo, propRepo)
	alertSvc  := services.NewAlertService(alertRepo)
	searchSvc := services.NewSearchService(searchRepo)

	// ── Handlers ──────────────────────────────────────────────────────────
	authH   := handlers.NewAuthHandler(authSvc)
	wsH     := handlers.NewWSHandler(appws.H)
	userH   := handlers.NewUserHandler(userSvc)
	propH   := handlers.NewPropertyHandler(propSvc)
	rentalH := handlers.NewRentalHandler(rentalSvc)
	walletH := handlers.NewWalletHandler(walletSvc)
	txH     := handlers.NewTransactionHandler(txSvc)
	reviewH := handlers.NewReviewHandler(reviewSvc)
	msgH    := handlers.NewMessageHandler(msgSvc)
	favH    := handlers.NewFavoriteHandler(favSvc)
	alertH  := handlers.NewAlertHandler(alertSvc)
	searchH := handlers.NewSearchHandler(searchSvc)

	// ── Routes (paths unchanged) ───────────────────────────────────────────
	api := r.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
		auth.GET("/me", middleware.Auth(), authH.GetMe)
	}

	r.GET("/ws", wsH.ServeWS)

	users := api.Group("/users", middleware.Auth())
	{
		users.PUT("/profile", userH.UpdateProfile)
		users.POST("/profile-picture", userH.UploadProfilePicture)
		users.PUT("/password", userH.ChangePassword)
		users.PATCH("/push-token", userH.UpdatePushToken)
		users.DELETE("/profile", userH.RequestDeleteAccount)
	}

	props := api.Group("/properties")
	{
		props.GET("", middleware.OptionalAuth(), propH.List)
		props.GET("/featured", propH.GetFeatured)
		props.GET("/:id", middleware.OptionalAuth(), propH.Get)
		props.POST("", middleware.Auth(), propH.Create)
		props.PUT("/:id", middleware.Auth(), propH.Update)
		props.DELETE("/:id", middleware.Auth(), propH.Delete)
		props.PATCH("/:id/availability", middleware.Auth(), propH.ToggleAvailability)
		props.GET("/my/listings", middleware.Auth(), propH.MyProperties)
		props.POST("/:id/images", middleware.Auth(), propH.AddImage)
		props.GET("/:id/reviews", reviewH.GetPropertyReviews)
		props.POST("/:id/reviews", middleware.Auth(), reviewH.CreateReview)
	}

	images := api.Group("/images", middleware.Auth())
	{
		images.DELETE("/:id", propH.DeleteImage)
		images.PATCH("/:id/main", propH.SetMainImage)
	}

	reviews := api.Group("/reviews", middleware.Auth())
	{
		reviews.DELETE("/:id", reviewH.DeleteReview)
	}

	favs := api.Group("/favorites", middleware.Auth())
	{
		favs.GET("", favH.GetFavorites)
		favs.POST("/:id/toggle", favH.ToggleFavorite)
	}

	rentals := api.Group("/rentals", middleware.Auth())
	{
		rentals.GET("", rentalH.GetMyRentals)
		rentals.POST("", rentalH.CreateRental)
		rentals.PATCH("/:id/status", rentalH.UpdateRentalStatus)
	}

	wallets := api.Group("/wallets", middleware.Auth())
	{
		wallets.GET("", walletH.GetWallets)
		wallets.POST("", walletH.CreateWallet)
		wallets.PUT("/:id", walletH.UpdateWallet)
		wallets.DELETE("/:id", walletH.DeleteWallet)
		wallets.PATCH("/:id/select", walletH.SelectWallet)
		wallets.POST("/:id/refresh-balance", walletH.RefreshWalletBalance)
	}

	txs := api.Group("/transactions", middleware.Auth())
	{
		txs.GET("", txH.GetTransactions)
		txs.POST("/pay-service", txH.PayService)
		txs.POST("/transfer", txH.TransferMoney)
	}

	alerts := api.Group("/alerts", middleware.Auth())
	{
		alerts.GET("", alertH.GetAlerts)
		alerts.POST("", alertH.CreateAlert)
		alerts.PUT("/:id", alertH.UpdateAlert)
		alerts.DELETE("/:id", alertH.DeleteAlert)
	}

	convs := api.Group("/conversations", middleware.Auth())
	{
		convs.GET("", msgH.GetConversations)
		convs.POST("", msgH.StartConversation)
		convs.GET("/:id", msgH.GetConversation)
		convs.POST("/:id/messages", msgH.SendMessage)
		convs.DELETE("/:id/messages/:msgID", msgH.DeleteMessage)
	}

	search := api.Group("/search")
	{
		search.GET("/suggestions", searchH.GetSearchSuggestions)
		search.POST("/history", middleware.OptionalAuth(), searchH.SaveSearchHistory)
		search.GET("/history", middleware.Auth(), searchH.GetSearchHistory)
		search.DELETE("/history", middleware.Auth(), searchH.DeleteSearchHistory)
	}
}
