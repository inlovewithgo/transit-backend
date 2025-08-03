package routes

import (
    "github.com/gofiber/fiber/v2"
    "github.com/inlovewithgo/transit-backend/main/handlers/api/basic"
    // "github.com/inlovewithgo/transit-backend/main/middlewares"
)

func SetupRoutes(app *fiber.App) {
    api := app.Group("/api/v1")

    health := api.Group("/health")
    {
        health.Get("/", handlers.BasicHealthCheck)
        health.Get("/detailed", handlers.DetailedHealthCheck)
        health.Get("/ready", handlers.ReadinessCheck)
        health.Get("/live", handlers.LivenessCheck)
    }

   // protected := api.Group("/", middleware.AuthMiddleware())
   // {
   // }

    // Root health check (alternative endpoint)
    app.Get("/health", handlers.BasicHealthCheck)
    app.Get("/", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "message": "Transit Backend API",
            "status":  "running",
            "version": "1.0.0",
        })
    })
}
