package seeder

import (
	"log"

	"dinacom-11.0-backend/models/entity"

	"gorm.io/gorm"
)

var achievements = []entity.Achievement{
	// Reporting Milestones
	{ID: "first_report", Name: "Langkah Pertama", Description: "Report pertama berhasil dibuat", BadgeURL: "/badges/first_report.webp", Category: "milestone", XPReward: 50},
	{ID: "road_warrior_5", Name: "Pemantau Jalan", Description: "Submit 5 laporan kerusakan jalan", BadgeURL: "/badges/road_warrior_5.webp", Category: "milestone", XPReward: 100},
	{ID: "road_warrior_25", Name: "Pengawal Jalanan", Description: "Submit 25 laporan kerusakan jalan", BadgeURL: "/badges/road_warrior_25.webp", Category: "milestone", XPReward: 250},
	{ID: "road_warrior_100", Name: "Pahlawan Infrastruktur", Description: "Submit 100 laporan kerusakan jalan", BadgeURL: "/badges/road_warrior_100.webp", Category: "milestone", XPReward: 500},

	// Quality Achievements
	{ID: "high_impact", Name: "Dampak Tinggi", Description: "Laporan dengan total score lebih dari 80", BadgeURL: "/badges/high_impact.webp", Category: "quality", XPReward: 150},
	{ID: "critical_finder", Name: "Penemu Kritis", Description: "Melaporkan kerusakan dengan kelas 'Very Poor'", BadgeURL: "/badges/critical_finder.webp", Category: "quality", XPReward: 200},
	{ID: "verified_contributor", Name: "Kontributor Terverifikasi", Description: "Laporan pertama yang diverifikasi admin", BadgeURL: "/badges/verified_contributor.webp", Category: "quality", XPReward: 100},
	{ID: "perfect_10", Name: "Sepuluh Sempurna", Description: "10 laporan yang telah diverifikasi", BadgeURL: "/badges/perfect_10.webp", Category: "quality", XPReward: 300},

	// Location Explorer
	{ID: "street_explorer", Name: "Penjelajah Jalan", Description: "Melaporkan kerusakan di 5 jalan berbeda", BadgeURL: "/badges/street_explorer.webp", Category: "explorer", XPReward: 150},
	{ID: "city_guardian", Name: "Penjaga Kota", Description: "Melaporkan kerusakan di 10 jalan berbeda", BadgeURL: "/badges/city_guardian.webp", Category: "explorer", XPReward: 300},

	// Streak Achievements
	{ID: "weekly_active", Name: "Aktif Mingguan", Description: "Membuat laporan 7 hari berturut-turut", BadgeURL: "/badges/weekly_active.webp", Category: "streak", XPReward: 200},
	{ID: "early_bird", Name: "Si Pagi", Description: "Membuat laporan sebelum jam 7 pagi", BadgeURL: "/badges/early_bird.webp", Category: "streak", XPReward: 75},
}

func SeedAchievements(db *gorm.DB) {
	for _, achievement := range achievements {
		var existing entity.Achievement
		result := db.Where("id = ?", achievement.ID).First(&existing)
		if result.RowsAffected == 0 {
			if err := db.Create(&achievement).Error; err != nil {
				log.Printf("Failed to seed achievement %s: %v", achievement.ID, err)
			} else {
				log.Printf("Seeded achievement: %s", achievement.Name)
			}
		}
	}
	log.Println("Achievement seeding completed")
}
