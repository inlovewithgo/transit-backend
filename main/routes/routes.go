package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/inlovewithgo/transit-backend/main/config"
	handlers "github.com/inlovewithgo/transit-backend/main/handlers/api/basic"
	authHandlers "github.com/inlovewithgo/transit-backend/main/handlers/auth"
	waitlistHandlers "github.com/inlovewithgo/transit-backend/main/handlers/waitlist"
	"github.com/inlovewithgo/transit-backend/main/middlewares"
	"github.com/inlovewithgo/transit-backend/main/repo/postgres"
	"github.com/inlovewithgo/transit-backend/main/service"
)

func SetupRoutes(app *fiber.App) {
	db := config.GetDB()

	// Repositories
	userRepo := postgres.NewUserRepository(db)
	waitlistRepo := postgres.NewWaitlistRepository(db)

	// Services
	mailService := service.NewMailService()
	authService := service.NewAuthService(userRepo, mailService)
	waitlistService := service.NewWaitlistService(waitlistRepo, mailService)

	// Handlers
	authHandler := authHandlers.NewAuthHandler(authService)
	waitlistHandler := waitlistHandlers.NewWaitlistHandler(waitlistService)

	// Rate limiter
	rateLimiter := middlewares.NewRateLimiter()

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

	// Waitlist routes with rate limiting
	waitlist := api.Group("/waitlist")
	{
		waitlist.Post("/", rateLimiter.WaitlistRateLimit(), waitlistHandler.JoinWaitlist)
		waitlist.Get("/stats", waitlistHandler.GetWaitlistStats)
	}

	protected := api.Group("/", middlewares.AuthMiddleware())
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
