package provider

import "dinacom-11.0-backend/repositories"

type RepositoriesProvider interface {
	ProvideConnectionRepository() repositories.ConnectionRepository
}

type repositoriesProvider struct {
	connectionRepository repositories.ConnectionRepository
}

func NewRepositoriesProvider(cfg ConfigProvider) RepositoriesProvider {
	connectionRepository := repositories.NewConnectionRepository(cfg.ProvideDatabaseConfig().GetInstance())
	return &repositoriesProvider{connectionRepository: connectionRepository}
}

func (rp *repositoriesProvider) ProvideConnectionRepository() repositories.ConnectionRepository {
	return rp.connectionRepository
}
