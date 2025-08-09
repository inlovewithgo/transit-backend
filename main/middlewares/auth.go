package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/inlovewithgo/transit-backend/main/models"
	"github.com/inlovewithgo/transit-backend/main/utils"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Missing authorization header",
			})
		}

		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Invalid authorization header format",
			})
		}

		tokenString := authHeader[7:]

		claims, err := utils.ValidateAccessToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Invalid or expired token",
			})
		}

		c.Locals("userID", claims.UserID)
		c.Locals("userEmail", claims.Email)

		return c.Next()
	}
}
