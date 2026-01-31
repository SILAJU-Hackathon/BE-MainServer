package provider

import "dinacom-11.0-backend/repositories"

type RepositoriesProvider interface {
	ProvideUserRepository() repositories.UserRepository
	ProvideReportRepository() repositories.ReportRepository
	ProvideAchievementRepository() repositories.AchievementRepository
	ProvideNotificationRepository() repositories.NotificationRepository
}

type repositoriesProvider struct {
	userRepository         repositories.UserRepository
	reportRepository       repositories.ReportRepository
	achievementRepository  repositories.AchievementRepository
	notificationRepository repositories.NotificationRepository
}

func NewRepositoriesProvider(cfg ConfigProvider) RepositoriesProvider {
	userRepository := repositories.NewUserRepository(cfg.ProvideDatabaseConfig().GetInstance())
	reportRepository := repositories.NewReportRepository(cfg.ProvideDatabaseConfig().GetInstance())
	achievementRepository := repositories.NewAchievementRepository(cfg.ProvideDatabaseConfig().GetInstance())
	notificationRepository := repositories.NewNotificationRepository(cfg.ProvideDatabaseConfig().GetInstance())
	return &repositoriesProvider{
		userRepository:         userRepository,
		reportRepository:       reportRepository,
		achievementRepository:  achievementRepository,
		notificationRepository: notificationRepository,
	}
}

func (rp *repositoriesProvider) ProvideUserRepository() repositories.UserRepository {
	return rp.userRepository
}

func (rp *repositoriesProvider) ProvideReportRepository() repositories.ReportRepository {
	return rp.reportRepository
}

func (rp *repositoriesProvider) ProvideAchievementRepository() repositories.AchievementRepository {
	return rp.achievementRepository
}

func (rp *repositoriesProvider) ProvideNotificationRepository() repositories.NotificationRepository {
	return rp.notificationRepository
}
