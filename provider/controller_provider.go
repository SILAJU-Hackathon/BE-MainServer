package provider

import "dinacom-11.0-backend/controllers"

type ControllerProvider interface {
	ProvideConnectionController() controllers.ConnectionController
}

type controllerProvider struct {
	connectionController controllers.ConnectionController
}

func NewControllerProvider(servicesProvider ServicesProvider) ControllerProvider {

	connectionController := controllers.NewConnectionController(servicesProvider.ProvideConnectionService())
	return &controllerProvider{
		connectionController: connectionController,
	}
}

func (c *controllerProvider) ProvideConnectionController() controllers.ConnectionController {
	return c.connectionController
}
