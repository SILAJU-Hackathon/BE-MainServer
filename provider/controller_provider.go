package provider

import "dinacom-11.0-backend/controllers"

type ControllerProvider interface {
	ProvideAuthController() controllers.AuthController
	ProvideReportController() controllers.ReportController
	ProvideAchievementController() controllers.AchievementController
}

type controllerProvider struct {
	authController        controllers.AuthController
	reportController      controllers.ReportController
	achievementController controllers.AchievementController
}

func NewControllerProvider(servicesProvider ServicesProvider) ControllerProvider {
	authController := controllers.NewAuthController(servicesProvider.ProvideAuthService())
	reportController := controllers.NewReportController(servicesProvider.ProvideReportService())
	achievementController := controllers.NewAchievementController(servicesProvider.ProvideAchievementService())
	return &controllerProvider{
		authController:        authController,
		reportController:      reportController,
		achievementController: achievementController,
	}
}

func (c *controllerProvider) ProvideAuthController() controllers.AuthController {
	return c.authController
}

func (c *controllerProvider) ProvideReportController() controllers.ReportController {
	return c.reportController
}

func (c *controllerProvider) ProvideAchievementController() controllers.AchievementController {
	return c.achievementController
}
