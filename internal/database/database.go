package database

import (
	"log"

	"github.com/badersalis/gidana_backend/internal/config"
	"github.com/badersalis/gidana_backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	var err error
	cfg := config.App

	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	if cfg.DatabaseURL != "" {
		DB, err = gorm.Open(postgres.Open(cfg.DatabaseURL), gormCfg)
	} else {
		DB, err = gorm.Open(sqlite.Open(cfg.DBPath), gormCfg)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")
	migrate()
}

func preMigrate() {
	type colType struct {
		DataType string
	}

	// If properties.id is still UUID the whole schema was UUID-based.
	// Drop everything and let AutoMigrate rebuild with integer PKs.
	var pkCol colType
	DB.Raw(`
		SELECT data_type
		FROM information_schema.columns
		WHERE table_schema = CURRENT_SCHEMA()
		  AND table_name   = 'properties'
		  AND column_name  = 'id'
	`).Scan(&pkCol)
	if pkCol.DataType == "uuid" {
		log.Println("preMigrate: UUID primary keys detected — dropping all tables for schema reset")
		tables := []string{
			"messages", "conversations", "search_histories",
			"transactions", "wallets", "alerts", "favorites",
			"reviews", "rentals", "property_images", "properties",
			"deleted_accounts", "users",
		}
		for _, t := range tables {
			if err := DB.Exec("DROP TABLE IF EXISTS \"" + t + "\" CASCADE").Error; err != nil {
				log.Fatalf("preMigrate: failed to drop table %s: %v", t, err)
			}
		}
		log.Println("preMigrate: all tables dropped — AutoMigrate will recreate them")
		return
	}

	// Narrower fix: property_images.property_id was uuid but properties.id is now bigint.
	var fkCol colType
	DB.Raw(`
		SELECT data_type
		FROM information_schema.columns
		WHERE table_schema = CURRENT_SCHEMA()
		  AND table_name   = 'property_images'
		  AND column_name  = 'property_id'
	`).Scan(&fkCol)
	if fkCol.DataType == "uuid" {
		if err := DB.Exec(`ALTER TABLE property_images DROP COLUMN property_id`).Error; err != nil {
			log.Fatalf("preMigrate: failed to drop uuid property_id: %v", err)
		}
		log.Println("preMigrate: dropped uuid property_id from property_images")
	}
}

func migrate() {
	preMigrate()
	err := DB.AutoMigrate(
		&models.User{},
		&models.Property{},
		&models.PropertyImage{},
		&models.Rental{},
		&models.Review{},
		&models.Favorite{},
		&models.Alert{},
		&models.SearchHistory{},
		&models.Conversation{},
		&models.Message{},
		&models.DeletedAccount{},
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Database migration completed")
}
