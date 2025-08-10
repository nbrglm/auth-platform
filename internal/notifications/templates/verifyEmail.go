package templates

import (
	"html/template"
	"path"

	"github.com/nbrglm/auth-platform/config"
)

func newVerifyEmailTemplate() (*EmailTemplate, error) {
	if config.Notifications.Email.TemplatesDir != nil {
		subjectTemplate, err := template.ParseFiles(path.Join(*config.Notifications.Email.TemplatesDir, "data/verify-email/subject.txt"))
		if err != nil {
			return nil, err
		}
		htmlTemplate, err := template.ParseFiles(path.Join(*config.Notifications.Email.TemplatesDir, "data/verify-email/body.html"))
		if err != nil {
			return nil, err
		}
		plainTextTemplate, err := template.ParseFiles(path.Join(*config.Notifications.Email.TemplatesDir, "data/verify-email/plain-text.txt"))
		if err != nil {
			return nil, err
		}
		return &EmailTemplate{
			TemplateName:  "verify_email",
			Subject:       subjectTemplate,
			HTMLBody:      htmlTemplate,
			PlainTextBody: plainTextTemplate,
		}, nil
	}

	subjectTemplate, err := template.ParseFS(templateFs, "data/verify-email/subject.txt")
	if err != nil {
		return nil, err
	}
	htmlTemplate, err := template.ParseFS(templateFs, "data/verify-email/body.html")
	if err != nil {
		return nil, err
	}
	plainTextTemplate, err := template.ParseFS(templateFs, "data/verify-email/plain-text.txt")
	if err != nil {
		return nil, err
	}
	return &EmailTemplate{
		TemplateName:  "verify_email",
		Subject:       subjectTemplate,
		HTMLBody:      htmlTemplate,
		PlainTextBody: plainTextTemplate,
	}, nil
}
