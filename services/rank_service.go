package services

import (
	"math"

	"dinacom-11.0-backend/models/dto"
	"dinacom-11.0-backend/repositories"

	"github.com/google/uuid"
)

const (
	XP_BASE       = 100.0
	XP_MULTIPLIER = 1.5
	MAX_LEVEL     = 10
)

// Rank definitions
var ranks = []struct {
	MinLevel int
	Name     string
}{
	{1, "Bronze Reporter"},
	{3, "Silver Reporter"},
	{5, "Gold Reporter"},
	{7, "Platinum Reporter"},
	{9, "Diamond Reporter"},
}

type RankService interface {
	GetUserRank(userID uuid.UUID) (*dto.UserRankResponse, error)
	AddXPToUser(userID uuid.UUID, xp int) (*dto.LevelUpResponse, error)
	CalculateLevel(totalXP int) int
	GetXPForLevel(level int) int
	GetCumulativeXPForLevel(level int) int
	GetLeaderboard(limit int) ([]dto.LeaderboardEntry, error)
}

type rankService struct {
	userRepo repositories.UserRepository
}

func NewRankService(userRepo repositories.UserRepository) RankService {
	return &rankService{userRepo: userRepo}
}

// GetXPForLevel calculates XP needed to reach next level from current level
// Formula: XP_required = Constant × Multiplier^Level
func (s *rankService) GetXPForLevel(level int) int {
	if level <= 1 {
		return 0
	}
	return int(XP_BASE * math.Pow(XP_MULTIPLIER, float64(level-1)))
}

// GetCumulativeXPForLevel calculates total XP needed to reach a level
func (s *rankService) GetCumulativeXPForLevel(level int) int {
	total := 0
	for l := 2; l <= level; l++ {
		total += s.GetXPForLevel(l)
	}
	return total
}

// CalculateLevel determines user level based on total XP
func (s *rankService) CalculateLevel(totalXP int) int {
	for level := MAX_LEVEL; level >= 1; level-- {
		if totalXP >= s.GetCumulativeXPForLevel(level) {
			return level
		}
	}
	return 1
}

// getRank returns rank info based on level
func (s *rankService) getRank(level int) (int, string) {
	rank := 1
	name := ranks[0].Name

	for i, r := range ranks {
		if level >= r.MinLevel {
			rank = i + 1
			name = r.Name
		}
	}
	return rank, name
}

func (s *rankService) GetUserRank(userID uuid.UUID) (*dto.UserRankResponse, error) {
	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}

	level := s.CalculateLevel(user.TotalXP)
	currentLevelXP := s.GetCumulativeXPForLevel(level)
	nextLevelXP := s.GetCumulativeXPForLevel(level + 1)
	xpInCurrentLevel := user.TotalXP - currentLevelXP
	xpNeededForNext := nextLevelXP - currentLevelXP

	var progress float64 = 100.0
	nextLevel := level
	if level < MAX_LEVEL {
		progress = float64(xpInCurrentLevel) / float64(xpNeededForNext) * 100
		nextLevel = level + 1
	}

	rank, rankName := s.getRank(level)

	return &dto.UserRankResponse{
		Level:     level,
		Rank:      rank,
		RankName:  rankName,
		TotalXP:   user.TotalXP,
		CurrentXP: xpInCurrentLevel,
		XPToNext:  xpNeededForNext - xpInCurrentLevel,
		Progress:  math.Round(progress*100) / 100,
		NextLevel: nextLevel,
	}, nil
}

func (s *rankService) AddXPToUser(userID uuid.UUID, xp int) (*dto.LevelUpResponse, error) {
	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}

	oldLevel := s.CalculateLevel(user.TotalXP)
	oldRank, _ := s.getRank(oldLevel)

	user.TotalXP += xp
	newLevel := s.CalculateLevel(user.TotalXP)
	user.Level = newLevel

	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, err
	}

	newRank, newRankName := s.getRank(newLevel)
	rankChanged := newRank > oldRank

	response := &dto.LevelUpResponse{
		LeveledUp:   newLevel > oldLevel,
		OldLevel:    oldLevel,
		NewLevel:    newLevel,
		XPGained:    xp,
		RankChanged: rankChanged,
	}

	if rankChanged {
		response.NewRankName = newRankName
	}

	return response, nil
}

func (s *rankService) GetLeaderboard(limit int) ([]dto.LeaderboardEntry, error) {
	users, err := s.userRepo.GetTopUsersByXP(limit)
	if err != nil {
		return nil, err
	}

	var leaderboard []dto.LeaderboardEntry
	for i, user := range users {
		_, rankName := s.getRank(user.Level)

		leaderboard = append(leaderboard, dto.LeaderboardEntry{
			Rank:     i + 1,
			UserID:   user.ID.String(),
			Fullname: user.Fullname,
			TotalXP:  user.TotalXP,
			Level:    user.Level,
			RankName: rankName,
		})
	}

	return leaderboard, nil
}
