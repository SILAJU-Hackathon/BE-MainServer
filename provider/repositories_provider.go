package provider

import "dinacom-11.0-backend/repositories"

type RepositoriesProvider interface {
	ProvideConnectionRepository() repositories.ConnectionRepository
	ProvideUserRepository() repositories.UserRepository
	ProvideReportRepository() repositories.ReportRepository
}

type repositoriesProvider struct {
	connectionRepository repositories.ConnectionRepository
	userRepository       repositories.UserRepository
	reportRepository     repositories.ReportRepository
}

func NewRepositoriesProvider(cfg ConfigProvider) RepositoriesProvider {
	connectionRepository := repositories.NewConnectionRepository(cfg.ProvideDatabaseConfig().GetInstance())
	userRepository := repositories.NewUserRepository(cfg.ProvideDatabaseConfig().GetInstance())
	reportRepository := repositories.NewReportRepository(cfg.ProvideDatabaseConfig().GetInstance())
	return &repositoriesProvider{
		connectionRepository: connectionRepository,
		userRepository:       userRepository,
		reportRepository:     reportRepository,
	}
}

func (rp *repositoriesProvider) ProvideConnectionRepository() repositories.ConnectionRepository {
	return rp.connectionRepository
}

func (rp *repositoriesProvider) ProvideUserRepository() repositories.UserRepository {
	return rp.userRepository
}

func (rp *repositoriesProvider) ProvideReportRepository() repositories.ReportRepository {
	return rp.reportRepository
}
