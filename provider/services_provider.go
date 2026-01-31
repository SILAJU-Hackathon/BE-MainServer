package provider

import "dinacom-11.0-backend/services"

type ServicesProvider interface {
	ProvideConnectionService() services.ConnectionService
	ProvideAuthService() services.AuthService
	ProvideReportService() services.ReportService
}

type servicesProvider struct {
	connectionService services.ConnectionService
	authService       services.AuthService
	reportService     services.ReportService
}

func NewServicesProvider(repoProvider RepositoriesProvider, configProvider ConfigProvider) ServicesProvider {
	connectionService := services.NewConnectionService(repoProvider.ProvideConnectionRepository())
	authService := services.NewAuthService(repoProvider.ProvideUserRepository())
	reportService := services.NewReportService(repoProvider.ProvideReportRepository())
	return &servicesProvider{
		connectionService: connectionService,
		authService:       authService,
		reportService:     reportService,
	}
}

func (s *servicesProvider) ProvideConnectionService() services.ConnectionService {
	return s.connectionService
}

func (s *servicesProvider) ProvideAuthService() services.AuthService {
	return s.authService
}

func (s *servicesProvider) ProvideReportService() services.ReportService {
	return s.reportService
}
