package provider

import "dinacom-11.0-backend/controllers"

type ControllerProvider interface {
	ProvideConnectionController() controllers.ConnectionController
	ProvideAuthController() controllers.AuthController
}

type controllerProvider struct {
	connectionController controllers.ConnectionController
	authController       controllers.AuthController
}

func NewControllerProvider(servicesProvider ServicesProvider) ControllerProvider {

	connectionController := controllers.NewConnectionController(servicesProvider.ProvideConnectionService())
	authController := controllers.NewAuthController(servicesProvider.ProvideAuthService())
	return &controllerProvider{
		connectionController: connectionController,
		authController:       authController,
	}
}

func (c *controllerProvider) ProvideConnectionController() controllers.ConnectionController {
	return c.connectionController
}

func (c *controllerProvider) ProvideAuthController() controllers.AuthController {
	return c.authController
}
