package provider

import "dinacom-11.0-backend/controllers"

type ControllerProvider interface {
	ProvideAuthController() controllers.AuthController
	ProvideReportController() controllers.ReportController
	ProvideAchievementController() controllers.AchievementController
	ProvideRankController() controllers.RankController
}

type controllerProvider struct {
	authController        controllers.AuthController
	reportController      controllers.ReportController
	achievementController controllers.AchievementController
	rankController        controllers.RankController
}

func NewControllerProvider(servicesProvider ServicesProvider) ControllerProvider {
	authController := controllers.NewAuthController(servicesProvider.ProvideAuthService())
	reportController := controllers.NewReportController(servicesProvider.ProvideReportService())
	achievementController := controllers.NewAchievementController(servicesProvider.ProvideAchievementService())
	rankController := controllers.NewRankController(servicesProvider.ProvideRankService())
	return &controllerProvider{
		authController:        authController,
		reportController:      reportController,
		achievementController: achievementController,
		rankController:        rankController,
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

func (c *controllerProvider) ProvideRankController() controllers.RankController {
	return c.rankController
}
