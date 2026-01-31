package router

import (
	"dinacom-11.0-backend/provider"
)

func RunRouter(appProvider provider.AppProvider) {
	router, controller, config := appProvider.ProvideRouter(), appProvider.ProvideControllers(), appProvider.ProvideConfig()
	ConnectionRouter(router, controller)
	err := router.Run(config.ProvideEnvConfig().GetTCPAddress())
	if err != nil {
		panic(err)
	}
}
