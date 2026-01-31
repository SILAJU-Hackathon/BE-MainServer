package provider

import "dinacom-11.0-backend/services"

type ServicesProvider interface {
	ProvideAuthService() services.AuthService
	ProvideReportService() services.ReportService
}

type servicesProvider struct {
	authService   services.AuthService
	reportService services.ReportService
}

func NewServicesProvider(repoProvider RepositoriesProvider, configProvider ConfigProvider) ServicesProvider {
	authService := services.NewAuthService(repoProvider.ProvideUserRepository())
	reportService := services.NewReportService(repoProvider.ProvideReportRepository())
	return &servicesProvider{
		authService:   authService,
		reportService: reportService,
	}
}

func (s *servicesProvider) ProvideAuthService() services.AuthService {
	return s.authService
}

func (s *servicesProvider) ProvideReportService() services.ReportService {
	return s.reportService
}
