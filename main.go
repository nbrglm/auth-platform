package main

import (
	"github.com/nbrglm/auth-platform/cmd"
)

// @title NBRGLM Auth Platform
// @version 0.0.1
// @description This is the NBRGLM Auth Platform API documentation.
// @termsOfService https://nbrglm.com/auth-platform/terms

// @contact.name NBRGLM Support
// @contact.url https://nbrglm.com/support
// @contact.email contact@nbrglm.com

// @license.name Apache 2.0
// @license.url https://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3360
// @BasePath /

// @securityDefinitions.apiKey APIKeyAuth
// @in header
// @name X-NAP-API-Key

// @securityDefinitions.apikey SessionHeaderAuth
// @in header
// @name X-NAP-Session-Token

// @securityDefinitions.apiKey RefreshHeaderAuth
// @in header
// @name X-NAP-Refresh-Token

// @externalDocs.description NBRGLM Auth Platform Documentation
// @externalDocs.url https://nbrglm.com/auth-platform/docs
func main() {
	// Execute the root command.
	cmd.Exec()
}
