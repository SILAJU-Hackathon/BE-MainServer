package services

import (
	"time"

	"dinacom-11.0-backend/models/dto"
	"dinacom-11.0-backend/repositories"

	"github.com/google/uuid"
)

type AchievementService interface {
	GetUserAchievements(userID uuid.UUID) (*dto.AchievementListResponse, error)
	GetUnlockedAchievements(userID uuid.UUID) (*dto.AchievementListResponse, error)
	CheckAndUnlockAchievements(userID uuid.UUID) ([]dto.NewAchievementResponse, error)
}

type achievementService struct {
	achievementRepo repositories.AchievementRepository
}

func NewAchievementService(achievementRepo repositories.AchievementRepository) AchievementService {
	return &achievementService{achievementRepo: achievementRepo}
}

func (s *achievementService) GetUserAchievements(userID uuid.UUID) (*dto.AchievementListResponse, error) {
	allAchievements, err := s.achievementRepo.GetAllAchievements()
	if err != nil {
		return nil, err
	}

	userAchievements, err := s.achievementRepo.GetUserAchievements(userID)
	if err != nil {
		return nil, err
	}

	unlockedMap := make(map[string]time.Time)
	for _, ua := range userAchievements {
		unlockedMap[ua.AchievementID] = ua.UnlockedAt
	}

	var achievements []dto.AchievementResponse
	unlockedCount := 0
	for _, a := range allAchievements {
		unlocked := false
		unlockedAt := ""
		if t, ok := unlockedMap[a.ID]; ok {
			unlocked = true
			unlockedAt = t.Format(time.RFC3339)
			unlockedCount++
		}
		achievements = append(achievements, dto.AchievementResponse{
			ID:          a.ID,
			Name:        a.Name,
			Description: a.Description,
			BadgeURL:    a.BadgeURL,
			Category:    a.Category,
			Unlocked:    unlocked,
			UnlockedAt:  unlockedAt,
		})
	}

	return &dto.AchievementListResponse{
		Achievements:  achievements,
		TotalCount:    len(allAchievements),
		UnlockedCount: unlockedCount,
	}, nil
}

func (s *achievementService) GetUnlockedAchievements(userID uuid.UUID) (*dto.AchievementListResponse, error) {
	userAchievements, err := s.achievementRepo.GetUserAchievements(userID)
	if err != nil {
		return nil, err
	}

	var achievements []dto.AchievementResponse
	for _, ua := range userAchievements {
		achievements = append(achievements, dto.AchievementResponse{
			ID:          ua.Achievement.ID,
			Name:        ua.Achievement.Name,
			Description: ua.Achievement.Description,
			BadgeURL:    ua.Achievement.BadgeURL,
			Category:    ua.Achievement.Category,
			Unlocked:    true,
			UnlockedAt:  ua.UnlockedAt.Format(time.RFC3339),
		})
	}

	return &dto.AchievementListResponse{
		Achievements:  achievements,
		TotalCount:    len(achievements),
		UnlockedCount: len(achievements),
	}, nil
}

func (s *achievementService) CheckAndUnlockAchievements(userID uuid.UUID) ([]dto.NewAchievementResponse, error) {
	var newAchievements []dto.NewAchievementResponse

	// Get all achievements for badge info
	allAchievements, _ := s.achievementRepo.GetAllAchievements()
	achievementMap := make(map[string]struct{ Name, Description, BadgeURL string })
	for _, a := range allAchievements {
		achievementMap[a.ID] = struct{ Name, Description, BadgeURL string }{a.Name, a.Description, a.BadgeURL}
	}

	// Helper function
	unlock := func(achievementID string) {
		if err := s.achievementRepo.UnlockAchievement(userID, achievementID); err == nil {
			if info, ok := achievementMap[achievementID]; ok {
				newAchievements = append(newAchievements, dto.NewAchievementResponse{
					AchievementID: achievementID,
					Name:          info.Name,
					Description:   info.Description,
					BadgeURL:      info.BadgeURL,
				})
			}
		}
	}

	// Check report count milestones
	reportCount, _ := s.achievementRepo.GetUserReportCount(userID)
	if reportCount >= 1 {
		if has, _ := s.achievementRepo.HasAchievement(userID, "first_report"); !has {
			unlock("first_report")
		}
	}
	if reportCount >= 5 {
		if has, _ := s.achievementRepo.HasAchievement(userID, "road_warrior_5"); !has {
			unlock("road_warrior_5")
		}
	}
	if reportCount >= 25 {
		if has, _ := s.achievementRepo.HasAchievement(userID, "road_warrior_25"); !has {
			unlock("road_warrior_25")
		}
	}
	if reportCount >= 100 {
		if has, _ := s.achievementRepo.HasAchievement(userID, "road_warrior_100"); !has {
			unlock("road_warrior_100")
		}
	}

	// Check quality achievements
	if hasHigh, _ := s.achievementRepo.HasHighImpactReport(userID); hasHigh {
		if has, _ := s.achievementRepo.HasAchievement(userID, "high_impact"); !has {
			unlock("high_impact")
		}
	}

	if hasCritical, _ := s.achievementRepo.HasCriticalReport(userID); hasCritical {
		if has, _ := s.achievementRepo.HasAchievement(userID, "critical_finder"); !has {
			unlock("critical_finder")
		}
	}

	verifiedCount, _ := s.achievementRepo.GetUserVerifiedReportCount(userID)
	if verifiedCount >= 1 {
		if has, _ := s.achievementRepo.HasAchievement(userID, "verified_contributor"); !has {
			unlock("verified_contributor")
		}
	}
	if verifiedCount >= 10 {
		if has, _ := s.achievementRepo.HasAchievement(userID, "perfect_10"); !has {
			unlock("perfect_10")
		}
	}

	// Check location explorer
	uniqueRoads, _ := s.achievementRepo.GetUserUniqueRoadCount(userID)
	if uniqueRoads >= 5 {
		if has, _ := s.achievementRepo.HasAchievement(userID, "street_explorer"); !has {
			unlock("street_explorer")
		}
	}
	if uniqueRoads >= 10 {
		if has, _ := s.achievementRepo.HasAchievement(userID, "city_guardian"); !has {
			unlock("city_guardian")
		}
	}

	// Check streak achievements
	dates, _ := s.achievementRepo.GetUserReportDatesLast7Days(userID)
	if len(dates) >= 7 {
		if has, _ := s.achievementRepo.HasAchievement(userID, "weekly_active"); !has {
			unlock("weekly_active")
		}
	}

	if hasEarly, _ := s.achievementRepo.HasEarlyMorningReport(userID); hasEarly {
		if has, _ := s.achievementRepo.HasAchievement(userID, "early_bird"); !has {
			unlock("early_bird")
		}
	}

	return newAchievements, nil
}
