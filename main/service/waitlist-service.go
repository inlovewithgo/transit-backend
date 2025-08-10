package service

import (
	"errors"
	"strings"

	"github.com/inlovewithgo/transit-backend/main/models"
	repo "github.com/inlovewithgo/transit-backend/main/repo/interface"
	"github.com/inlovewithgo/transit-backend/pkg/logger"
	"gorm.io/gorm"
)

type WaitlistService struct {
	waitlistRepo repo.WaitlistRepository
	mailService  *MailService
}

func NewWaitlistService(waitlistRepo repo.WaitlistRepository, mailService *MailService) *WaitlistService {
	return &WaitlistService{
		waitlistRepo: waitlistRepo,
		mailService:  mailService,
	}
}

func (ws *WaitlistService) JoinWaitlist(email, ipAddress, userAgent string) (*models.WaitlistResponse, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	if !isValidEmail(email) {
		return &models.WaitlistResponse{
			Message: "Invalid email format",
			Status:  "error",
		}, errors.New("invalid email format")
	}

	existingUser, err := ws.waitlistRepo.GetByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Log.Error("Error checking existing waitlist user: %v", err)
		return &models.WaitlistResponse{
			Message: "Internal server error",
			Status:  "error",
		}, err
	}

	if existingUser != nil {
		logger.Log.Info("User %s already exists in waitlist", email)
		return &models.WaitlistResponse{
			Message: "You're already on the waitlist! We'll notify you when we launch.",
			Status:  "already_registered",
			Email:   email,
		}, nil
	}

	waitlistEntry := &models.Waitlist{
		Email:     email,
		Status:    "pending",
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}

	contactID, err := ws.mailService.AddToWaitlistAudience(email)
	if err != nil {
		logger.Log.Error("Failed to add %s to Resend audience: %v", email, err)
		waitlistEntry.Status = "failed"
	} else {
		waitlistEntry.ResendContactID = contactID
		waitlistEntry.Status = "subscribed"
	}

	err = ws.waitlistRepo.Create(waitlistEntry)
	if err != nil {
		logger.Log.Error("Failed to save waitlist entry: %v", err)
		return &models.WaitlistResponse{
			Message: "Failed to join waitlist. Please try again.",
			Status:  "error",
		}, err
	}

	if waitlistEntry.Status == "subscribed" {
		go func() {
			if err := ws.mailService.SendWaitlistConfirmationEmail(email); err != nil {
				logger.Log.Error("Failed to send confirmation email to %s: %v", email, err)
			}
		}()
	}

	response := &models.WaitlistResponse{
		Email: email,
	}

	if waitlistEntry.Status == "subscribed" {
		response.Message = "Successfully joined the waitlist! Check your email for confirmation."
		response.Status = "success"
	} else {
		response.Message = "Added to waitlist, but failed to subscribe to email updates."
		response.Status = "partial_success"
	}

	logger.Log.Info("User %s successfully joined waitlist with status: %s", email, waitlistEntry.Status)
	return response, nil
}

func (ws *WaitlistService) GetWaitlistStats() (map[string]interface{}, error) {
	count, err := ws.waitlistRepo.GetCount()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_subscribers": count,
		"message":           "Waitlist statistics",
	}, nil
}

func isValidEmail(email string) bool {
	if len(email) < 5 || len(email) > 254 {
		return false
	}

	atIndex := strings.LastIndex(email, "@")
	if atIndex < 1 || atIndex == len(email)-1 {
		return false
	}

	domain := email[atIndex+1:]
	return strings.Contains(domain, ".")
}
