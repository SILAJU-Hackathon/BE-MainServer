package provider

import "dinacom-11.0-backend/services"

type ServicesProvider interface {
	ProvideAuthService() services.AuthService
	ProvideReportService() services.ReportService
	ProvideAchievementService() services.AchievementService
	ProvideRankService() services.RankService
}

type servicesProvider struct {
	authService        services.AuthService
	reportService      services.ReportService
	achievementService services.AchievementService
	rankService        services.RankService
}

func NewServicesProvider(repoProvider RepositoriesProvider, configProvider ConfigProvider) ServicesProvider {
	authService := services.NewAuthService(repoProvider.ProvideUserRepository())
	rankService := services.NewRankService(repoProvider.ProvideUserRepository())
	reportService := services.NewReportService(repoProvider.ProvideReportRepository(), repoProvider.ProvideUserRepository(), rankService)
	achievementService := services.NewAchievementService(repoProvider.ProvideAchievementRepository(), rankService)
	return &servicesProvider{
		authService:        authService,
		reportService:      reportService,
		achievementService: achievementService,
		rankService:        rankService,
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

func (s *servicesProvider) ProvideRankService() services.RankService {
	return s.rankService
}
