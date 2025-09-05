package templates

func newAdminLoginTemplate() (*EmailTemplate, error) {
	htmlTmplSubPath := "templs/AdminLogin/body.html"
	plainTextTmplSubPath := "templs/AdminLogin/plain-text.txt"
	subjectTmplSubPath := "templs/AdminLogin/subject.txt"

	subjectTemplate, htmlTemplate, plainTextTemplate, err := findAndParseTemplates(htmlTmplSubPath, plainTextTmplSubPath, subjectTmplSubPath)
	if err != nil {
		return nil, err
	}

	return &EmailTemplate{
		TemplateName:  "AdminLogin",
		Subject:       subjectTemplate,
		HTMLBody:      htmlTemplate,
		PlainTextBody: plainTextTemplate,
	}, nil
}
