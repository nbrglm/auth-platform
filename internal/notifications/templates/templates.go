// Package templates provides the templates for notifications.
//
// This package contains the templates used for sending notifications like SignUP, Password Reset, and Email Verification, etc.
package templates

import (
	"bytes"
	"embed"
	"html/template"
	"strings"
	"time"

	"github.com/nbrglm/auth-platform/config"
)

//go:embed data
var templateFs embed.FS

type TemplateData struct {
	AppName     string
	UserName    string
	UserEmail   string
	ActionURL   string
	ExpiresAt   time.Time
	IPAddress   string
	UserAgent   string
	Location    string
	SupportURL  string
	CompanyName string

	// NBRGLMBranding is a flag to indicate whether the email should include NBRGLM branding.
	//
	// Usually set to config.NBRGLMBranding, which is a boolean value indicating whether to include NBRGLM branding.
	NBRGLMBranding bool
}

// EmailTemplate represents the structure of an email template.
//
// It contains fields for the template name, subject, HTML body, and plain text body.
// The TemplateName is used for identification purposes, while the Subject, HTMLBody, and PlainTextBody
// are used to define the content of the email. Only these three fields are rendered for the user's email.
type EmailTemplate struct {
	// TemplateName is the name of the template, used for identification.
	TemplateName  string
	Subject       *template.Template
	HTMLBody      *template.Template
	PlainTextBody *template.Template
}

// RenderedEmailTemplate represents the rendered email template.
//
// It contains the template name, subject, HTML body, and plain text body.
// This struct is used to hold the final rendered content after processing the email template with data.
// It is useful for sending the email with the actual content filled in.
type RenderedEmailTemplate struct {
	TemplateName  string
	Subject       string
	HTMLBody      string
	PlainTextBody string
}

type MessageTemplate struct {
}

// The following variables store the parsed email and sms templates.
//
// Any template that needs to be used for sending notifications should be defined here.
var (
	// VerifyEmailTemplate is the template used for verifying email addresses.
	VerifyEmailTemplate *EmailTemplate
)

// Must be called to parse all email templates at application startup.
// This function initializes the email templates used for notifications.
func ParseEmailTemplates() (err error) {
	VerifyEmailTemplate, err = newVerifyEmailTemplate()
	if err != nil {
		return err
	}
	return nil
}

// Must be called to parse all message templates at application startup.
// This function initializes the sms templates used for notifications.
func ParseMessageTemplates() (err error) {
	return nil
}

// RenderEmailTemplate renders the email template with the provided data.
// It takes a TemplateData struct and an EmailTemplate struct as input,
// and returns a RenderedEmailTemplate struct with the rendered content.
// The function is expected to replace placeholders in the email template with actual data from TemplateData.
// If the rendering fails, it returns an error and the original template.
func RenderEmailTemplate(data TemplateData, tmpl EmailTemplate) (*RenderedEmailTemplate, error) {
	var htmlBody, plainTextBody, subject bytes.Buffer
	data.NBRGLMBranding = config.NBRGLMBranding // Ensure NBRGLMBranding is set based on the configuration
	if err := tmpl.HTMLBody.ExecuteTemplate(&htmlBody, "VerifyEmailHTML", data); err != nil {
		return nil, err
	}
	if err := tmpl.PlainTextBody.ExecuteTemplate(&plainTextBody, "VerifyEmailText", data); err != nil {
		return nil, err
	}
	if err := tmpl.Subject.ExecuteTemplate(&subject, "VerifyEmailSubject", data); err != nil {
		return nil, err
	}
	return &RenderedEmailTemplate{
		TemplateName:  tmpl.TemplateName,
		Subject:       strings.TrimSpace(subject.String()),
		HTMLBody:      strings.TrimSpace(htmlBody.String()),
		PlainTextBody: strings.TrimSpace(plainTextBody.String()),
	}, nil
}
