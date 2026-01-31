package provider

import "dinacom-11.0-backend/controllers"

type ControllerProvider interface {
	ProvideAuthController() controllers.AuthController
	ProvideReportController() controllers.ReportController
}

type controllerProvider struct {
	authController   controllers.AuthController
	reportController controllers.ReportController
}

func NewControllerProvider(servicesProvider ServicesProvider) ControllerProvider {
	authController := controllers.NewAuthController(servicesProvider.ProvideAuthService())
	reportController := controllers.NewReportController(servicesProvider.ProvideReportService())
	return &controllerProvider{
		authController:   authController,
		reportController: reportController,
	}
}

func (c *controllerProvider) ProvideAuthController() controllers.AuthController {
	return c.authController
}

func (c *controllerProvider) ProvideReportController() controllers.ReportController {
	return c.reportController
}
