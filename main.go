package main

import (
	"dinacom-11.0-backend/provider"
	"dinacom-11.0-backend/router"
)

// @title Dinacom 11.0 Backend API
// @version 1.0
// @description Backend service for Dinacom 11.0 Hackathon
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@dinacom.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	appProvider := provider.NewAppProvider()
	router.RunRouter(appProvider)
}
