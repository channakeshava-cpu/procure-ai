package models

import "time"

type RecommendationSession struct {
	ID                string    `json:"id" gorm:"primaryKey;size:32"`
	Category          string    `json:"category" gorm:"not null"`
	Quantity          int       `json:"quantity" gorm:"not null"`
	Budget            float64   `json:"budget"`
	MaxDeliveryDays   int       `json:"maxDeliveryDays"`
	PreferredCities   string    `json:"preferredCities,omitempty" gorm:"type:text;not null;default:''"`
	TopN              int       `json:"topN" gorm:"not null;default:5"`
	RecommendedVendor string    `json:"recommendedVendor,omitempty" gorm:"not null;default:''"`
	ShortlistSnapshot string    `json:"shortlistSnapshot,omitempty" gorm:"type:text;not null;default:''"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
