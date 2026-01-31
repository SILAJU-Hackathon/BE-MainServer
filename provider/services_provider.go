package provider

import "dinacom-11.0-backend/services"

type ServicesProvider interface {
	ProvideConnectionService() services.ConnectionService
	ProvideAuthService() services.AuthService
}

type servicesProvider struct {
	connectionService services.ConnectionService
	authService       services.AuthService
}

func NewServicesProvider(repoProvider RepositoriesProvider, configProvider ConfigProvider) ServicesProvider {
	connectionService := services.NewConnectionService(repoProvider.ProvideConnectionRepository())
	authService := services.NewAuthService(repoProvider.ProvideUserRepository())
	return &servicesProvider{
		connectionService: connectionService,
		authService:       authService,
	}
}

func (s *servicesProvider) ProvideConnectionService() services.ConnectionService {
	return s.connectionService
}

func (s *servicesProvider) ProvideAuthService() services.AuthService {
	return s.authService
}
