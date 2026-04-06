package db

import (
	"procure-ai/models"

	"gorm.io/gorm"
)

func AutoMigrate(database *gorm.DB) error {
	return database.AutoMigrate(
		&models.Vendor{},
		&models.Order{},
		&models.QR{},
		&models.RecommendationSession{},
	)
}
