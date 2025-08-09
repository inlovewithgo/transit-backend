package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/inlovewithgo/transit-backend/main/config"
	handlers "github.com/inlovewithgo/transit-backend/main/handlers/api/basic"
	authHandlers "github.com/inlovewithgo/transit-backend/main/handlers/auth"
	middleware "github.com/inlovewithgo/transit-backend/main/middlewares"
	"github.com/inlovewithgo/transit-backend/main/repo/postgres"
	"github.com/inlovewithgo/transit-backend/main/service"
)

func SetupRoutes(app *fiber.App) {
	db := config.GetDB()

	userRepo := postgres.NewUserRepository(db)

	mailService := service.NewMailService()
	authService := service.NewAuthService(userRepo, mailService)

	authHandler := authHandlers.NewAuthHandler(authService)

	api := app.Group("/api/v1")

	health := api.Group("/health")
	{
		health.Get("/", handlers.BasicHealthCheck)
		health.Get("/detailed", handlers.DetailedHealthCheck)
		health.Get("/ready", handlers.ReadinessCheck)
		health.Get("/live", handlers.LivenessCheck)
	}

	auth := api.Group("/auth")
	{
		auth.Post("/register", authHandler.Register)
		auth.Post("/login", authHandler.Login)
	}

	protected := api.Group("/", middleware.AuthMiddleware())
	{
		protected.Get("/profile", authHandler.GetProfile)
		protected.Post("/logout", authHandler.Logout)
	}

	app.Get("/health", handlers.BasicHealthCheck)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Transit Backend API",
			"status":  "running",
			"version": "1.0.0",
		})
	})
}
