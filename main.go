package main

import (
	"embed"

	"github.com/nbrglm/auth-platform/cmd"
)

//go:embed web
var webFS embed.FS

// @title NBRGLM Auth Platform
// @version 0.1.0
// @description This is the NBRGLM Auth Platform API documentation.
// @termsOfService https://nbrglm.com/auth-platform/terms

// @contact.name NBRGLM Support
// @contact.url https://nbrglm.com/support
// @contact.email contact@nbrglm.com

// @license.name Apache 2.0
// @license.url https://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey SessionHeaderAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey SessionCookieAuth
// @in cookie
// @name nap-session-tk

// @securityDefinitions.apikey CSRFTokenCookieAuth
// @in cookie
// @name nap-csrf-tk

// @securityDefinitions.apikey CSRFTokenHeaderAuth
// @in header
// @name X-CSRF-Token

// @securityDefinitions.apiKey RefreshCookieAuth
// @in cookie
// @name nap-refresh-tk

// @securityDefinitions.apiKey RefreshHeaderAuth
// @in header
// @name X-NAP-Refresh-Token

// @externalDocs.description NBRGLM Auth Platform Documentation
// @externalDocs.url https://nbrglm.com/auth-platform/docs
func main() {
	// Execute the root command.
	cmd.Exec(webFS)
}
