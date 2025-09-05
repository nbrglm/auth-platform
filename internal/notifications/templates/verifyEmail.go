package templates

func newVerifyEmailTemplate() (*EmailTemplate, error) {
	htmlTmplSubPath := "templs/VerifyEmail/body.html"
	plainTextTmplSubPath := "templs/VerifyEmail/plain-text.txt"
	subjectTmplSubPath := "templs/VerifyEmail/subject.txt"

	subjectTemplate, htmlTemplate, plainTextTemplate, err := findAndParseTemplates(htmlTmplSubPath, plainTextTmplSubPath, subjectTmplSubPath)
	if err != nil {
		return nil, err
	}

	return &EmailTemplate{
		TemplateName:  "VerifyEmail",
		Subject:       subjectTemplate,
		HTMLBody:      htmlTemplate,
		PlainTextBody: plainTextTemplate,
	}, nil
}
