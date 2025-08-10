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
    htmlContent := `
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>Welcome to Transit Waitlist</title>
        </head>
        <body style="margin: 0; padding: 0; background-color: #f4f4f4; font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;">
            <div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
                <!-- Header -->
                <div style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); padding: 40px 30px; text-align: center;">
                    <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 600;">
                        You're on the waitlist!
                    </h1>
                    <p style="color: #e8e8e8; margin: 10px 0 0 0; font-size: 16px;">
                        Welcome to the Transit Backend early community
                    </p>
                </div>

                <!-- Main Content -->
                <div style="padding: 40px 30px;">
                    <div style="text-align: center; margin-bottom: 30px;">
                        <div style="width: 80px; height: 80px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); border-radius: 50%; margin: 0 auto 20px; display: flex; align-items: center; justify-content: center; font-size: 32px;">
                            🎉
                        </div>
                    </div>

                    <p style="color: #333333; font-size: 16px; line-height: 1.6; margin-bottom: 20px;">
                        Hello there!
                    </p>
                    
                    <p style="color: #333333; font-size: 16px; line-height: 1.6; margin-bottom: 30px;">
                        Thank you for joining the <strong>Transit waitlist</strong>! We're thrilled to have you as part of our early community. You'll be among the first to experience our innovative transit solutions.
                    </p>

                    <!-- What's Next Section -->
                    <div style="background: linear-gradient(135deg, #f8f9ff 0%, #f0f2ff 100%); padding: 25px; border-radius: 12px; margin: 30px 0; border-left: 4px solid #667eea;">
                        <h3 style="color: #333333; margin: 0 0 15px 0; font-size: 18px; font-weight: 600;">
                            What's next?
                        </h3>
                        <div style="color: #555555; font-size: 14px; line-height: 1.6;">
                            <div style="display: flex; align-items: flex-start; margin-bottom: 12px;">
                                <span style="color: #667eea; font-weight: bold; margin-right: 10px;"></span>
                                <span>Keep an eye on your inbox for exclusive updates and announcements</span>
                            </div>
                            <div style="display: flex; align-items: flex-start; margin-bottom: 12px;">
                                <span style="color: #667eea; font-weight: bold; margin-right: 10px;"></span>
                                <span>Be the first to access new features and beta releases</span>
                            </div>
                            <div style="display: flex; align-items: flex-start; margin-bottom: 12px;">
                                <span style="color: #667eea; font-weight: bold; margin-right: 10px;"></span>
                                <span>Get priority support when we officially launch</span>
                            </div>
                            <div style="display: flex; align-items: flex-start;">
                                <span style="color: #667eea; font-weight: bold; margin-right: 10px;"></span>
                                <span>Enjoy exclusive early-bird benefits and special offers</span>
                            </div>
                        </div>
                    </div>

                    <!-- Call to Action -->
                    <div style="text-align: center; margin: 30px 0;">
                        <div style="background-color: #f8f9fa; padding: 20px; border-radius: 8px; border: 2px dashed #dee2e6;">
                            <p style="margin: 0; color: #6c757d; font-size: 14px;">
                                💡 <strong>Pro tip:</strong> Add our email to your contacts to ensure you never miss important updates!
                            </p>
                        </div>
                    </div>

                    <p style="color: #333333; font-size: 16px; line-height: 1.6; text-align: center;">
                        Stay tuned for exciting updates! 🚀
                    </p>
                </div>

                <!-- Footer -->
                <div style="background-color: #f8f9fa; padding: 25px 30px; text-align: center; border-top: 1px solid #dee2e6;">
                    <p style="color: #666666; font-size: 14px; margin: 0 0 10px 0;">
                        Best regards,<br>
                        <strong style="color: #333333;">The Transit Team</strong>
                    </p>
                    <div style="margin-top: 15px; padding-top: 15px; border-top: 1px solid #dee2e6;">
                        <p style="color: #999999; font-size: 12px; margin: 0;">
                            This email was sent to you because you joined our waitlist. 
                            We're excited to have you on board! 🎉
                        </p>
                    </div>
                </div>
            </div>
        </body>
        </html>
    `

    params := &resend.SendEmailRequest{
        From:    "noreply@yssh.dev",
        To:      []string{email},
        Subject: "Welcome to Transit Waitlist!",
        Html:    htmlContent,
    }

    _, err := ms.client.Emails.Send(params)
    if err != nil {
        logger.Log.Error("Failed to send waitlist confirmation email to %s: %v", email, err)
        return err
    }

    logger.Log.Info("Waitlist confirmation email sent successfully to %s", email)
    return nil
}
