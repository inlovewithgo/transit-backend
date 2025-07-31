package db

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/inlovewithgo/transit-backend/main/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DbName       string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  int
}

type GormLoggerImpl struct {}

func LoadConfig() (*Config, error) {
	port, err := strconv.Atoi(utils.GetENV("DB_PORT", "5432"))
	if err != nil {
		return nil, err
	}
	maxOpen, err := strconv.Atoi(utils.GetENV("DB_MAX_OPEN_CONNS", "100"))
	if err != nil {
		return nil, err
	}
	maxIdle, err := strconv.Atoi(utils.GetENV("DB_MAX_IDLE_CONNS", "10"))
	if err != nil {
		return nil, err
	}
	maxLifetime, err := strconv.Atoi(utils.GetENV("DB_MAX_LIFETIME", "3600"))
	if err != nil {
		return nil, err
	}

	return &Config{
		Host:         utils.GetENV("DB_HOST", "localhost"),
		Port:         port,
		Username:     utils.GetENV("DB_USER", "user"),
		Password:     utils.GetENV("DB_PASSWORD", "password"),
		DbName:       utils.GetENV("DB_NAME", "dbname"),
		SSLMode:      utils.GetENV("DB_SSL_MODE", "disable"),
		MaxOpenConns: maxOpen,
		MaxIdleConns: maxIdle,
		MaxLifetime:  maxLifetime,
	}, nil
}



func NewGormDB(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DbName, cfg.SSLMode)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default,
		})

	if err != nil {
		fmt.Printf("failed to connect to database: %v\n", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("failed to get fetch sql.DB from GORM: %v\n", err)
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifetime))

	ctx, cancel := context.WithTimeout(context.Background(), 4 * time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		fmt.Printf("PostGres ping failed: %v\n", err)
		return nil, err
	}

	fmt.Println("PostGres connection established successfully")

	return db, nil
}

func CloseGormDB(db *gorm.DB) {
	if db == nil {
		return
	}

	sqlDB, err := db.DB()
	if err == nil {
		if err := sqlDB.Close(); err != nil {
			fmt.Printf("failed to close database connection: %v\n", err)
		} else {
			fmt.Println("Database connection closed successfully")
		}
	}
}