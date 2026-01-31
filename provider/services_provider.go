package provider

import "dinacom-11.0-backend/services"

type ServicesProvider interface {
	ProvideAuthService() services.AuthService
	ProvideReportService() services.ReportService
	ProvideAchievementService() services.AchievementService
}

type servicesProvider struct {
	authService        services.AuthService
	reportService      services.ReportService
	achievementService services.AchievementService
}

func NewServicesProvider(repoProvider RepositoriesProvider, configProvider ConfigProvider) ServicesProvider {
	authService := services.NewAuthService(repoProvider.ProvideUserRepository())
	reportService := services.NewReportService(repoProvider.ProvideReportRepository(), repoProvider.ProvideUserRepository())
	achievementService := services.NewAchievementService(repoProvider.ProvideAchievementRepository())
	return &servicesProvider{
		authService:        authService,
		reportService:      reportService,
		achievementService: achievementService,
	}
}

func (s *servicesProvider) ProvideAuthService() services.AuthService {
	return s.authService
}

func (s *servicesProvider) ProvideReportService() services.ReportService {
	return s.reportService
}

func (s *servicesProvider) ProvideAchievementService() services.AchievementService {
	return s.achievementService
}
