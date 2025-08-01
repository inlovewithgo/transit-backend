package config

import (
	"sync"

	"github.com/inlovewithgo/transit-backend/pkg/db"
	"github.com/inlovewithgo/transit-backend/pkg/logger"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
	once sync.Once
)

func InitDatabase() {
	once.Do(func() {
		cfg, err := db.LoadConfig()
		if err != nil {
			logger.Log.Fatal("Unable to load Database configuration: ", err)
		}

		inst, err := db.NewGormDB(cfg)
		if err != nil {
			logger.Log.Fatal("Unable to connect to Database: ", err)
		}

		DB = inst
	})
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
	}
}