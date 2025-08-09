package db

import (
	"context"
	"fmt"
	"log"
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

func LoadConfig() (*Config, error) {
	port, err := strconv.Atoi(utils.GetENV("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	maxOpen, err := strconv.Atoi(utils.GetENV("DB_MAX_OPEN_CONNS", "100"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_OPEN_CONNS: %w", err)
	}

	maxIdle, err := strconv.Atoi(utils.GetENV("DB_MAX_IDLE_CONNS", "10"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_IDLE_CONNS: %w", err)
	}

	maxLifetime, err := strconv.Atoi(utils.GetENV("DB_MAX_LIFETIME", "3600"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_LIFETIME: %w", err)
	}

	user := utils.GetENV("DB_USER", "postgres")
	password := utils.GetENV("DB_PASSWORD", "")
	dbName := utils.GetENV("DB_NAME", "transit-wallet")
	sslMode := utils.GetENV("DB_SSL_MODE", "disable")

	if user == "" || password == "" || dbName == "" {
		return nil, fmt.Errorf("missing required database environment variables")
	}

	return &Config{
		Host:         utils.GetENV("DB_HOST", "localhost"),
		Port:         port,
		Username:     user,
		Password:     password,
		DbName:       dbName,
		SSLMode:      sslMode,
		MaxOpenConns: maxOpen,
		MaxIdleConns: maxIdle,
		MaxLifetime:  maxLifetime,
	}, nil
}

func NewGormDB(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DbName, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from GORM: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifetime) * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	log.Println("PostgreSQL connection established successfully")
	return db, nil
}

func CloseGormDB(db *gorm.DB) {
	if db == nil {
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("failed to get sql.DB: %v\n", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Printf("failed to close database connection: %v\n", err)
	} else {
		log.Println("Database connection closed successfully")
	}
}
