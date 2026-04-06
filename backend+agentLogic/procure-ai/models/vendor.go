package models

import "time"

type Vendor struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	Name             string    `json:"name" binding:"required" gorm:"uniqueIndex;not null"`
	Price            float64   `json:"price" binding:"required,gt=0" gorm:"not null"`
	Trust            float64   `json:"trust" binding:"required,gte=0,lte=5" gorm:"not null"`
	DeliveryDays     int       `json:"deliveryDays" binding:"gte=0" gorm:"not null;default:0"`
	Stock            int       `json:"stock" binding:"gte=0" gorm:"not null;default:0"`
	MinOrderQty      int       `json:"minOrderQty" binding:"gte=0" gorm:"not null;default:1"`
	Location         string    `json:"location" gorm:"not null;default:''"`
	PaymentTerms     string    `json:"paymentTerms" gorm:"not null;default:''"`
	ReliabilityScore float64   `json:"reliabilityScore" binding:"gte=0,lte=100" gorm:"not null;default:0"`
	Category         string    `json:"category" gorm:"not null;default:''"`
	Orders           []Order   `json:"-" gorm:"foreignKey:VendorID"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type SelectVendorRequest struct {
	Vendors []Vendor `json:"vendors" binding:"required"`
}

type VendorSelectionResponse struct {
	RecommendedVendor Vendor         `json:"recommendedVendor"`
	Score             float64        `json:"score"`
	Reason            string         `json:"reason"`
	ScoredVendors     []ScoredVendor `json:"scoredVendors"`
}

type ScoredVendor struct {
	Vendor Vendor  `json:"vendor"`
	Score  float64 `json:"score"`
}
