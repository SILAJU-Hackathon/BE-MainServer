package provider

import "dinacom-11.0-backend/repositories"

type RepositoriesProvider interface {
	ProvideUserRepository() repositories.UserRepository
	ProvideReportRepository() repositories.ReportRepository
}

type repositoriesProvider struct {
	userRepository   repositories.UserRepository
	reportRepository repositories.ReportRepository
}

func NewRepositoriesProvider(cfg ConfigProvider) RepositoriesProvider {
	userRepository := repositories.NewUserRepository(cfg.ProvideDatabaseConfig().GetInstance())
	reportRepository := repositories.NewReportRepository(cfg.ProvideDatabaseConfig().GetInstance())
	return &repositoriesProvider{
		userRepository:   userRepository,
		reportRepository: reportRepository,
	}
}

func (rp *repositoriesProvider) ProvideUserRepository() repositories.UserRepository {
	return rp.userRepository
}

func (rp *repositoriesProvider) ProvideReportRepository() repositories.ReportRepository {
	return rp.reportRepository
}
