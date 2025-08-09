package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/inlovewithgo/transit-backend/main/models"
	"github.com/inlovewithgo/transit-backend/main/service"
	"github.com/inlovewithgo/transit-backend/main/utils"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid request format",
			Message: "Please provide valid JSON data",
		})
	}

	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Missing required fields",
			Message: "Email, password, first name, and last name are required",
		})
	}

	if len(req.Password) < 6 {
		return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Password too short",
			Message: "Password must be at least 6 characters long",
		})
	}

	response, err := h.authService.Register(&req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Registration failed",
			Message: err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(response)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid request format",
			Message: "Please provide valid JSON data",
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Missing required fields",
			Message: "Email and password are required",
		})
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "Login failed",
			Message: err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(response)
}

func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userID, err := getUserIDFromToken(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Invalid or missing token",
		})
	}

	user, err := h.authService.GetUserProfile(userID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(models.ErrorResponse{
			Error:   "User not found",
			Message: err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"user": user,
	})
}

func getUserIDFromToken(c *fiber.Ctx) (uint, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "Missing authorization header")
	}

	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization header format")
	}

	tokenString := authHeader[7:]
	claims, err := utils.ValidateAccessToken(tokenString)
	if err != nil {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	return claims.UserID, nil
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Logout successful",
	})
}
