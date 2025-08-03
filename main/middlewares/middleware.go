package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/recover"
    "github.com/gofiber/fiber/v2/middleware/requestid"
)

func SetupMiddleware(app *fiber.App) {
    app.Use(recover.New())

    app.Use(requestid.New())

    app.Use(cors.New(cors.Config{
        AllowOrigins: "*",
        AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
        AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    }))

    app.Use(logger.New(logger.Config{
        Format: "[${time}] ${ip} ${method} ${path} - ${status} - ${latency}\n",
        TimeFormat: "2006-01-02 15:04:05",
    }))
}

func AuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        
        auth := c.Get("Authorization")
        if auth == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Authorization header required",
            })
        }

        return c.Next()
    }
}

func RateLimitMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {

		// Implement your rate limiting logic here

        return c.Next()
    }
}
