package service

import (
	"fmt"
	"os"

	"github.com/inlovewithgo/transit-backend/pkg/logger"
	"github.com/inlovewithgo/transit-backend/main/utils"
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

func (ms *MailService) AddToWaitlistAudience(email string) (string, error) {
	audienceID := utils.GetENV("RESEND_AUDIENCE_ID", "")
	if audienceID == "" {
		logger.Log.Error("RESEND_AUDIENCE_ID environment variable is required")
		return "", fmt.Errorf("RESEND_AUDIENCE_ID not configured")
	}

	contactParams := &resend.CreateContactRequest{
		Email:      email,
		AudienceId: audienceID,
	}

	contact, err := ms.client.Contacts.Create(contactParams)
	if err != nil {
		logger.Log.Error("Failed to add %s to waitlist audience: %v", email, err)
		return "", err
	}

	logger.Log.Info("Successfully added %s to waitlist audience with contact ID: %s", email, contact.Id)
	return contact.Id, nil
}

func (ms *MailService) SendWaitlistConfirmationEmail(email string) error {
	params := &resend.SendEmailRequest{
		From:    "noreply@yssh.dev",
		To:      []string{email},
		Subject: "Welcome to Transit Backend Waitlist!",
		Html: fmt.Sprintf(`
			<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
				<h1 style="color: #333;">You're on the wait-list! 🎉</h1>
				<p>Hello,</p>
				<p>Thank you for joining the Transit waitlist! We're excited to have you as part of our early community.</p>
				<p>You'll be among the first to know when we launch new features and updates.</p>
				<div style="background-color: #f8f9fa; padding: 20px; border-radius: 8px; margin: 20px 0;">
					<h3 style="color: #333; margin-top: 0;">What's next?</h3>
					<ul style="color: #666;">
						<li>Keep an eye on your inbox for exclusive updates</li>
						<li>Be the first to access new features</li>
						<li>Get priority support when we launch</li>
					</ul>
				</div>
				<p>Stay tuned for exciting updates!</p>
				<br>
				<p>Best regards,<br>The Transit Team</p>
			</div>
		`),
	}

	_, err := ms.client.Emails.Send(params)
	if err != nil {
		logger.Log.Error("Failed to send waitlist confirmation email to %s: %v", email, err)
		return err
	}

	logger.Log.Info("Waitlist confirmation email sent successfully to %s", email)
	return nil
}
