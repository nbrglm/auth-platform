// Package notifications provides functionality to manage and send notifications.
// It includes support for email and SMS notifications, with configurations for each type.
//
// Email and SMS senders are implemented as interfaces, allowing for different implementations.
// Every implementation must have a config object, which is passed during initialization.
package notifications

import (
	"context"
	"fmt"
	"time"

	"github.com/nbrglm/auth-platform/config"
	"github.com/nbrglm/auth-platform/internal/logging"
	"github.com/nbrglm/auth-platform/internal/notifications/templates"
	"go.uber.org/zap"
)

type EmailSenderInterface interface {
	SendEmail(to string, subject string, htmlContent, plainTextContent string) error // SendEmail sends an email to the specified recipient with the given subject and body.
}

type SMSSenderInterface interface {
	SendSMS(to string, message string) error // SendSMS sends an SMS to the specified recipient with the given message.
}

var EmailSender EmailSenderInterface // EmailSender is the global email sender instance.
var SMSSender SMSSenderInterface     // SMSSender is the global SMS sender instance.

// EmailEnabled indicates whether email notifications are enabled.
var EmailEnabled = false

// SMSEnabled indicates whether SMS notifications are enabled.
var SMSEnabled = false

// InitEmail initializes the email sender based on the configuration.
// If the email sender is not configured, it logs a warning and skips initialization.
//
// Sets EmailEnabled to true if the email sender is configured.
func InitEmail() {
	if config.Notifications.Email == nil {
		logging.Logger.Warn("Email notifications are not configured, skipping email sender initialization... this will result in an error every time an email is sent")
		return
	}

	switch config.Notifications.Email.Provider {
	case "smtp":
		EmailSender = NewSMTPEmailSender(config.Notifications.Email.SMTP.Host, config.Notifications.Email.SMTP.Port, config.Notifications.Email.SMTP.FromAddress, config.Notifications.Email.SMTP.Password)
		EmailEnabled = true
	case "sendgrid":
		EmailSender = NewSendGridEmailSender(config.Notifications.Email.SendGrid.APIKey, config.Notifications.Email.SendGrid.FromAddress, config.Notifications.Email.SendGrid.FromName)
		EmailEnabled = true
	case "ses":
		EmailSender = NewSESEmailSender(config.Notifications.Email.SES.Region, config.Notifications.Email.SES.AccessKeyID, config.Notifications.Email.SES.SecretAccessKey, config.Notifications.Email.SES.FromAddress, config.Notifications.Email.SES.FromName)
		EmailEnabled = true
	default:
		EmailEnabled = false
		logging.Logger.Warn("Unknown email provider, skipping email sender initialization", zap.String("provider", config.Notifications.Email.Provider))
	}
}

func InitSMS() {
	if config.Notifications.SMS == nil {
		logging.Logger.Warn("SMS notifications are not configured, skipping SMS sender initialization... this will result in an error every time an SMS is sent")
		return
	}
}

var ErrEmailSenderNotSet = fmt.Errorf("email sender is not set, please set it using the config file! notifications.email.provider and the respective provider config")

type SendWelcomeEmailParams struct {
	User struct {
		Email     string
		FirstName *string
		LastName  *string
	}
	VerificationToken string
	ExpiresAt         time.Time
}

// SendWelcomeEmail sends a welcome email to the specified recipient.
// It uses the global EmailSender instance to send the email.
// The email also includes a link to verify the email address.
func SendWelcomeEmail(ctx context.Context, params SendWelcomeEmailParams) error {
	verificationUrl := fmt.Sprintf("%s/auth/verify-email?token=%s", config.Public.GetBaseURL(), params.VerificationToken)
	rendered, err := templates.RenderEmailTemplate(templates.TemplateData{
		AppName:     config.Branding.AppName,
		UserName:    getUserName(params.User.FirstName, params.User.LastName),
		UserEmail:   params.User.Email,
		ActionURL:   verificationUrl,
		ExpiresAt:   params.ExpiresAt,
		CompanyName: config.Branding.CompanyNameShort,
		SupportURL:  config.Branding.SupportURL,
	}, *templates.VerifyEmailTemplate)
	if err != nil {
		return err
	}

	logging.Logger.Debug("Sending welcome email", zap.String("to", params.User.Email), zap.String("subject", rendered.Subject), zap.String("html_body", rendered.HTMLBody), zap.String("plain_text_body", rendered.PlainTextBody))

	if EmailSender == nil {
		return ErrEmailSenderNotSet
	}

	err = EmailSender.SendEmail(params.User.Email, rendered.Subject, rendered.HTMLBody, rendered.PlainTextBody)
	if err != nil {
		return err
	}

	return nil
}

func getUserName(firstName, lastName *string) string {
	if firstName == nil && lastName == nil {
		return "User"
	}
	if firstName != nil && lastName != nil {
		return *firstName + " " + *lastName
	}
	if firstName != nil {
		return *firstName
	}
	return *lastName
}
