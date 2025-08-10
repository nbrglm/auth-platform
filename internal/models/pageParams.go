package models

import "github.com/nbrglm/auth-platform/config"

type CommonPageParams struct {
	Title            string
	CompanyName      string
	CompanyNameShort string
	SupportURL       string
	AppName          string
	CSRFTokenName    string
	CSRFTokenValue   string
	NBRGLMBranding   bool
	Multitenancy     bool
	PageError        *PageError
}

func NewCommonPageParams(title string, csrfToken string) CommonPageParams {
	return CommonPageParams{
		Title:            title,
		CompanyName:      config.Branding.CompanyName,
		CompanyNameShort: config.Branding.CompanyNameShort,
		SupportURL:       config.Branding.SupportURL,
		AppName:          config.Branding.AppName,
		CSRFTokenName:    config.Security.CSRF.TokenName,
		CSRFTokenValue:   csrfToken,
		NBRGLMBranding:   config.NBRGLMBranding,
		Multitenancy:     config.Multitenancy.Enable,
		PageError:        nil,
	}
}

type PageError struct {
	Body             string
	ReturnTo         string
	ReturnButtonText string
	HelpURL          string
	Error            bool
}

func NewPageError(title, body, returnTo, returnButtonText string) CommonPageParams {
	params := NewCommonPageParams(title, "")
	params.PageError = &PageError{
		Body:             body,
		ReturnTo:         returnTo,
		ReturnButtonText: returnButtonText,
		HelpURL:          config.Branding.SupportURL,
		Error:            true,
	}
	return params
}
