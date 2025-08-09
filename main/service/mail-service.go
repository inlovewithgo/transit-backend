package service

import (
	"fmt"
	"os"

	"github.com/inlovewithgo/transit-backend/pkg/logger"
	"github.com/resend/resend-go/v2"
)

type MailService struct {
	client *resend.Client
}

func NewMailService() *MailService {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		logger.Log.Fatal("RESEND_API_KEY environment variable is required")
	}

	client := resend.NewClient(apiKey)
	return &MailService{
		client: client,
	}
}

func (ms *MailService) SendWelcomeEmail(email, firstName, lastName string) error {
	params := &resend.SendEmailRequest{
		From:    "noreply@yssh.dev",
		To:      []string{email},
		Subject: "Welcome to Transit Backend!",
		Html: fmt.Sprintf(`
			<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
				<h1 style="color: #333;">Welcome to Transit Backend!</h1>
				<p>Hello %s %s,</p>
				<p>Thank you for registering with us. Your account has been successfully created!</p>
				<p>You can now start using our services.</p>
				<br>
				<p>Best regards,<br>The Transit Backend Team</p>
			</div>
		`, firstName, lastName),
	}

	_, err := ms.client.Emails.Send(params)
	if err != nil {
		logger.Log.Error("Failed to send welcome email to %s: %v", email, err)
		return err
	}

	logger.Log.Info("Welcome email sent successfully to %s", email)
	return nil
}

func (ms *MailService) SendLoginNotification(email, firstName, lastName string) error {
	params := &resend.SendEmailRequest{
		From:    "noreply@yssh.dev",
		To:      []string{email},
		Subject: "Successful Login - Transit Backend",
		Html: fmt.Sprintf(`
			<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
				<h1 style="color: #333;">Login Successful</h1>
				<p>Hello %s %s,</p>
				<p>You have successfully logged into your Transit Backend account.</p>
				<p>If this wasn't you, please contact our support team immediately.</p>
				<br>
				<p>Best regards,<br>The Transit Backend Team</p>
			</div>
		`, firstName, lastName),
	}

	_, err := ms.client.Emails.Send(params)
	if err != nil {
		logger.Log.Error("Failed to send login notification email to %s: %v", email, err)
		return err
	}

	logger.Log.Info("Login notification email sent successfully to %s", email)
	return nil
}
