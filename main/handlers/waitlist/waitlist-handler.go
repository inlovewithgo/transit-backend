package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/inlovewithgo/transit-backend/main/models"
	"github.com/inlovewithgo/transit-backend/main/service"
	"github.com/inlovewithgo/transit-backend/pkg/logger"
)

type WaitlistHandler struct {
	waitlistService *service.WaitlistService
}

func NewWaitlistHandler(waitlistService *service.WaitlistService) *WaitlistHandler {
	return &WaitlistHandler{
		waitlistService: waitlistService,
	}
}

// JoinWaitlist handles POST /api/v1/waitlist
func (wh *WaitlistHandler) JoinWaitlist(c *fiber.Ctx) error {
	var req models.WaitlistRequest

	if err := c.BodyParser(&req); err != nil {
		logger.Log.Error("Failed to parse waitlist request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request format",
		})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Email is required",
		})
	}

	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")

	response, err := wh.waitlistService.JoinWaitlist(req.Email, ipAddress, userAgent)
	if err != nil {
		if response != nil && response.Status == "error" {
			statusCode := fiber.StatusBadRequest
			if response.Message == "Internal server error" {
				statusCode = fiber.StatusInternalServerError
			}
			return c.Status(statusCode).JSON(fiber.Map{
				"error":   true,
				"message": response.Message,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to join waitlist",
		})
	}

	statusCode := fiber.StatusOK
	switch response.Status {
	case "success":
		statusCode = fiber.StatusCreated
	case "already_registered":
		statusCode = fiber.StatusOK
	case "partial_success":
		statusCode = fiber.StatusCreated
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"error":   false,
		"message": response.Message,
		"status":  response.Status,
		"data": fiber.Map{
			"email": response.Email,
		},
	})
}

func (wh *WaitlistHandler) GetWaitlistStats(c *fiber.Ctx) error {
	stats, err := wh.waitlistService.GetWaitlistStats()
	if err != nil {
		logger.Log.Error("Failed to get waitlist stats: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to retrieve waitlist statistics",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"data":  stats,
	})
}
