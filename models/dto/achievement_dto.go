package dto

type AchievementResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	BadgeURL    string `json:"badge_url"`
	Category    string `json:"category"`
	Unlocked    bool   `json:"unlocked"`
	UnlockedAt  string `json:"unlocked_at,omitempty"`
}

type AchievementListResponse struct {
	Achievements  []AchievementResponse `json:"achievements"`
	TotalCount    int                   `json:"total_count"`
	UnlockedCount int                   `json:"unlocked_count"`
}

type NewAchievementResponse struct {
	AchievementID string `json:"achievement_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	BadgeURL      string `json:"badge_url"`
}
