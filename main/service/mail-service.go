// TODO :
// 1. Revamp Emails HTML templates to something simple

package service

import (
    "fmt"
    "os"
	"time"

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
    htmlContent := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to Transit!</title>
</head>
<body style="margin: 0; padding: 40px 20px; background-color: #f5f5f5; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Arial, sans-serif;">
    <div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 10px rgba(0,0,0,0.1);">
        
        <!-- Header -->
        <div style="background-color: #000000; padding: 40px 30px; text-align: center;">
            <h1 style="color: #ffffff; margin: 0; font-size: 24px; font-weight: 600;">
                Welcome to Transit! 🚀
            </h1>
        </div>

        <!-- Content -->
        <div style="padding: 40px 30px; text-align: center;">
            <p style="color: #666666; font-size: 16px; margin: 0 0 20px 0;">
                Hi ` + firstName + ` ` + lastName + `! 👋
            </p>
            
            <p style="color: #333333; font-size: 16px; margin: 0 0 30px 0; line-height: 1.5;">
                Thank you for joining Transit! I'm excited to share insights on backend development, DevOps, and emerging technologies with you.
            </p>

            <p style="color: #333333; font-size: 16px; margin: 0 0 30px 0; line-height: 1.5;">
                You'll receive notifications whenever I publish new articles, tutorials, and technical deep-dives.
            </p>

            <a href="#" style="display: inline-block; background-color: #000000; color: #ffffff; text-decoration: none; padding: 12px 30px; border-radius: 6px; font-size: 16px; font-weight: 500; margin: 20px 0;">
                Start Exploring
            </a>

            <p style="color: #999999; font-size: 14px; margin: 30px 0 0 0;">
                If you no longer want to receive these emails, you can unsubscribe at any time.
            </p>
        </div>
    </div>
</body>
</html>`

    params := &resend.SendEmailRequest{
        From:    "noreply@yssh.dev",
        To:      []string{email},
        Subject: "🚀 Welcome to Transit!",
        Html:    htmlContent,
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
    currentTime := time.Now().Format("January 2, 2006 at 3:04 PM MST")
    
    htmlContent := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login Successful - Transit</title>
</head>
<body style="margin: 0; padding: 40px 20px; background-color: #f5f5f5; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Arial, sans-serif;">
    <div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 10px rgba(0,0,0,0.1);">
        
        <!-- Header -->
        <div style="background-color: #000000; padding: 40px 30px; text-align: center;">
            <h1 style="color: #ffffff; margin: 0; font-size: 24px; font-weight: 600;">
                Login Successful 🔐
            </h1>
        </div>

        <!-- Content -->
        <div style="padding: 40px 30px; text-align: center;">
            <p style="color: #666666; font-size: 16px; margin: 0 0 20px 0;">
                Welcome back, ` + firstName + ` ` + lastName + `! 👋
            </p>
            
            <p style="color: #333333; font-size: 16px; margin: 0 0 30px 0; line-height: 1.5;">
                You've successfully logged into your Transit account.
            </p>

            <!-- Session Info -->
            <div style="background-color: #f8f9fa; padding: 20px; border-radius: 6px; margin: 30px 0; text-align: left;">
                <div style="margin-bottom: 10px;">
                    <span style="color: #666666; font-size: 14px;">Account:</span>
                    <span style="color: #333333; font-size: 14px; float: right;">` + email + `</span>
                    <div style="clear: both;"></div>
                </div>
                <div style="margin-bottom: 10px;">
                    <span style="color: #666666; font-size: 14px;">Login Time:</span>
                    <span style="color: #333333; font-size: 14px; float: right;">` + currentTime + `</span>
                    <div style="clear: both;"></div>
                </div>
                <div>
                    <span style="color: #666666; font-size: 14px;">Status:</span>
                    <span style="color: #28a745; font-size: 14px; float: right;">✓ Secure</span>
                    <div style="clear: both;"></div>
                </div>
            </div>

            <p style="color: #666666; font-size: 14px; margin: 30px 0 0 0; line-height: 1.5;">
                If this wasn't you, please contact our support team immediately at <strong>security@yssh.dev</strong>
            </p>
        </div>
    </div>
</body>
</html>`

    params := &resend.SendEmailRequest{
        From:    "noreply@yssh.dev",
        To:      []string{email},
        Subject: "🔐 Login Successful - Transit",
        Html:    htmlContent,
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
    htmlContent := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to Transit Waitlist!</title>
</head>
<body style="margin: 0; padding: 40px 20px; background-color: #f5f5f5; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Arial, sans-serif;">
    <div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 10px rgba(0,0,0,0.1);">
        
        <!-- Header -->
        <div style="background-color: #000000; padding: 40px 30px; text-align: center;">
            <h1 style="color: #ffffff; margin: 0; font-size: 24px; font-weight: 600;">
                You're on the waitlist! 🎉
            </h1>
        </div>

        <!-- Content -->
        <div style="padding: 40px 30px; text-align: center;">
            <p style="color: #666666; font-size: 16px; margin: 0 0 20px 0;">
                Hello there! 👋
            </p>
            
            <p style="color: #333333; font-size: 16px; margin: 0 0 30px 0; line-height: 1.5;">
                Thank you for joining the <strong>Transit waitlist</strong>! We're thrilled to have you as part of our early community.
            </p>

            <p style="color: #333333; font-size: 16px; margin: 0 0 30px 0; line-height: 1.5;">
                You'll receive notifications whenever I publish new articles, tutorials, and technical deep-dives.
            </p>

            <!-- What's Next -->
            <div style="background-color: #f8f9fa; padding: 20px; border-radius: 6px; margin: 30px 0; text-align: left;">
                <h3 style="color: #333333; margin: 0 0 15px 0; font-size: 16px; font-weight: 600; text-align: center;">
                    What's next?
                </h3>
                <ul style="color: #666666; font-size: 14px; line-height: 1.6; margin: 0; padding-left: 20px;">
                    <li style="margin-bottom: 8px;">Keep an eye on your inbox for exclusive updates</li>
                    <li style="margin-bottom: 8px;">Be the first to access new features and beta releases</li>
                    <li style="margin-bottom: 8px;">Get priority support when we officially launch</li>
                    <li>Enjoy exclusive early-bird benefits and special offers</li>
                </ul>
            </div>

            <p style="color: #333333; font-size: 16px; margin: 30px 0 0 0; line-height: 1.5;">
                Stay tuned for exciting updates! 🚀
            </p>

            <p style="color: #999999; font-size: 14px; margin: 30px 0 0 0;">
                If you no longer want to receive these emails, you can unsubscribe at any time.
            </p>
        </div>
    </div>
</body>
</html>`

    params := &resend.SendEmailRequest{
        From:    "noreply@yssh.dev",
        To:      []string{email},
        Subject: "🎉 Welcome to Transit Waitlist!",
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