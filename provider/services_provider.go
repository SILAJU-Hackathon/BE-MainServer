package provider

import "dinacom-11.0-backend/services"

type ServicesProvider interface {
	ProvideConnectionService() services.ConnectionService
}

type servicesProvider struct {
	connectionService services.ConnectionService
}

func NewServicesProvider(repoProvider RepositoriesProvider, configProvider ConfigProvider) ServicesProvider {
	connectionService := services.NewConnectionService(repoProvider.ProvideConnectionRepository())
	return &servicesProvider{connectionService: connectionService}
}

func (s *servicesProvider) ProvideConnectionService() services.ConnectionService {
	return s.connectionService
}
