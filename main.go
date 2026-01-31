package main

import (
	"dinacom-11.0-backend/provider"
	"dinacom-11.0-backend/router"
)

func main() {
	appProvider := provider.NewAppProvider()
	router.RunRouter(appProvider)
}
