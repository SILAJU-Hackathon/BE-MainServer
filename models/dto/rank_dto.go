package dto

type UserRankResponse struct {
	Level     int     `json:"level"`
	Rank      int     `json:"rank"`
	RankName  string  `json:"rank_name"`
	TotalXP   int     `json:"total_xp"`
	CurrentXP int     `json:"current_xp"`
	XPToNext  int     `json:"xp_to_next"`
	Progress  float64 `json:"progress"`
	NextLevel int     `json:"next_level"`
}

type LevelUpResponse struct {
	LeveledUp   bool   `json:"leveled_up"`
	OldLevel    int    `json:"old_level"`
	NewLevel    int    `json:"new_level"`
	XPGained    int    `json:"xp_gained"`
	RankChanged bool   `json:"rank_changed"`
	NewRankName string `json:"new_rank_name,omitempty"`
}
