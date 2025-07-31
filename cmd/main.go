package main

import (
    "github.com/inlovewithgo/transit-backend/pkg/logger"
    "github.com/inlovewithgo/transit-backend/pkg/db"
)

func main() {
    log := logger.NewLogger()

    cfg, err := db.LoadConfig()
    if err != nil {
        log.Fatal("failed to load db config", zap.Error(err))
    }

    gormDB, err := db.NewGormDB(cfg, log)
    if err != nil {
        log.Fatal("failed to set up db", zap.Error(err))
    }
    defer db.CloseGormDB(gormDB, log)
}