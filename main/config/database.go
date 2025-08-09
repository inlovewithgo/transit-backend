package config

import (
	"sync"

	"github.com/inlovewithgo/transit-backend/main/models"
	"github.com/inlovewithgo/transit-backend/pkg/db"
	"github.com/inlovewithgo/transit-backend/pkg/logger"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	once sync.Once
)

func InitDatabase() {
	once.Do(func() {
		cfg, err := db.LoadConfig()
		if err != nil {
			logger.Log.Fatal("Unable to load Database configuration: %v", err)
		}

		inst, err := db.NewGormDB(cfg)
		if err != nil {
			logger.Log.Fatal("Unable to connect to Database: %v", err)
		}

		DB = inst

		// Run auto migrations
		if err := runMigrations(DB); err != nil {
			logger.Log.Fatal("Unable to run database migrations: %v", err)
		}
	})
}

func runMigrations(db *gorm.DB) error {
	// Auto-migrate all models
	err := db.AutoMigrate(
		&models.User{},
		// Add other models here as you create them
	)

	if err != nil {
		return err
	}

	logger.Log.Info("Database migrations completed successfully")
	return nil
}

func GetDB() *gorm.DB {
	if DB == nil {
		InitDatabase()
	}
	return DB
}

func ShutdownDatabase() {
	if DB != nil {
		db.CloseGormDB(DB)
		DB = nil
	}
}
