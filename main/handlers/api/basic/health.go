package handlers

import (
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/inlovewithgo/transit-backend/main/config"
)

type HealthResponse struct {
    Status    string    `json:"status"`
    Timestamp time.Time `json:"timestamp"`
    Database  string    `json:"database"`
    Version   string    `json:"version"`
    Uptime    string    `json:"uptime"`
}

type DetailedHealthResponse struct {
    Status    string                 `json:"status"`
    Timestamp time.Time              `json:"timestamp"`
    Services  map[string]ServiceInfo `json:"services"`
    System    SystemInfo             `json:"system"`
}

type ServiceInfo struct {
    Status      string `json:"status"`
    ResponseTime string `json:"response_time"`
    Message     string `json:"message,omitempty"`
}

type SystemInfo struct {
    Version string `json:"version"`
    Uptime  string `json:"uptime"`
}

var startTime = time.Now()

func BasicHealthCheck(c *fiber.Ctx) error {
    return c.Status(fiber.StatusOK).JSON(HealthResponse{
        Status:    "healthy",
        Timestamp: time.Now(),
        Database:  "connected",
        Version:   "1.0.0",
        Uptime:    time.Since(startTime).String(),
    })
}

func DetailedHealthCheck(c *fiber.Ctx) error {
    response := DetailedHealthResponse{
        Status:    "healthy",
        Timestamp: time.Now(),
        Services:  make(map[string]ServiceInfo),
        System: SystemInfo{
            Version: "1.0.0",
            Uptime:  time.Since(startTime).String(),
        },
    }

    dbStart := time.Now()
    db := config.GetDB()
    
    if db != nil {
        sqlDB, err := db.DB()
        if err != nil {
            response.Services["database"] = ServiceInfo{
                Status:       "unhealthy",
                ResponseTime: time.Since(dbStart).String(),
                Message:      "Failed to get database instance: " + err.Error(),
            }
            response.Status = "degraded"
        } else {
            err = sqlDB.Ping()
            if err != nil {
                response.Services["database"] = ServiceInfo{
                    Status:       "unhealthy",
                    ResponseTime: time.Since(dbStart).String(),
                    Message:      "Database ping failed: " + err.Error(),
                }
                response.Status = "degraded"
            } else {
                response.Services["database"] = ServiceInfo{
                    Status:       "healthy",
                    ResponseTime: time.Since(dbStart).String(),
                    Message:      "Database connection successful",
                }
            }
        }
    } else {
        response.Services["database"] = ServiceInfo{
            Status:       "unhealthy",
            ResponseTime: time.Since(dbStart).String(),
            Message:      "Database instance is nil",
        }
        response.Status = "degraded"
    }

    statusCode := fiber.StatusOK
    if response.Status == "degraded" {
        statusCode = fiber.StatusServiceUnavailable
    }

    return c.Status(statusCode).JSON(response)
}

func ReadinessCheck(c *fiber.Ctx) error {
    db := config.GetDB()
    if db == nil {
        return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
            "status": "not ready",
            "reason": "database not initialized",
        })
    }

    sqlDB, err := db.DB()
    if err != nil {
        return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
            "status": "not ready",
            "reason": "cannot access database",
        })
    }

    if err := sqlDB.Ping(); err != nil {
        return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
            "status": "not ready",
            "reason": "database connection failed",
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status": "ready",
    })
}

func LivenessCheck(c *fiber.Ctx) error {
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status": "alive",
        "timestamp": time.Now(),
    })
}
