package provider

import "dinacom-11.0-backend/repositories"

type RepositoriesProvider interface {
	ProvideConnectionRepository() repositories.ConnectionRepository
	ProvideUserRepository() repositories.UserRepository
}

type repositoriesProvider struct {
	connectionRepository repositories.ConnectionRepository
	userRepository       repositories.UserRepository
}

func NewRepositoriesProvider(cfg ConfigProvider) RepositoriesProvider {
	connectionRepository := repositories.NewConnectionRepository(cfg.ProvideDatabaseConfig().GetInstance())
	userRepository := repositories.NewUserRepository(cfg.ProvideDatabaseConfig().GetInstance())
	return &repositoriesProvider{
		connectionRepository: connectionRepository,
		userRepository:       userRepository,
	}
}

func (rp *repositoriesProvider) ProvideConnectionRepository() repositories.ConnectionRepository {
	return rp.connectionRepository
}

func (rp *repositoriesProvider) ProvideUserRepository() repositories.UserRepository {
	return rp.userRepository
}
