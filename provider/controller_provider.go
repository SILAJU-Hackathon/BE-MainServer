package provider

import "dinacom-11.0-backend/controllers"

type ControllerProvider interface {
	ProvideConnectionController() controllers.ConnectionController
	ProvideAuthController() controllers.AuthController
	ProvideReportController() controllers.ReportController
}

type controllerProvider struct {
	connectionController controllers.ConnectionController
	authController       controllers.AuthController
	reportController     controllers.ReportController
}

func NewControllerProvider(servicesProvider ServicesProvider) ControllerProvider {
	connectionController := controllers.NewConnectionController(servicesProvider.ProvideConnectionService())
	authController := controllers.NewAuthController(servicesProvider.ProvideAuthService())
	reportController := controllers.NewReportController(servicesProvider.ProvideReportService())
	return &controllerProvider{
		connectionController: connectionController,
		authController:       authController,
		reportController:     reportController,
	}
}

func (c *controllerProvider) ProvideConnectionController() controllers.ConnectionController {
	return c.connectionController
}

func (c *controllerProvider) ProvideAuthController() controllers.AuthController {
	return c.authController
}

func (c *controllerProvider) ProvideReportController() controllers.ReportController {
	return c.reportController
}
