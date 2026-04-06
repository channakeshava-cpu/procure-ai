package services

import (
	"errors"
	"fmt"
	"math"
	"sort"

	"procure-ai/models"

	"gorm.io/gorm"
)

type VendorService struct {
	db *gorm.DB
}

func NewVendorService(db *gorm.DB) *VendorService {
	return &VendorService{db: db}
}

func (s *VendorService) GetVendors() []models.Vendor {
	var vendors []models.Vendor
	s.db.Order("id ASC").Find(&vendors)
	return vendors
}

func (s *VendorService) GetVendorsByCategory(category string) []models.Vendor {
	var vendors []models.Vendor
	query := s.db.Order("id ASC")
	if category != "" {
		query = query.Where("category = ?", category)
	}
	query.Find(&vendors)
	return vendors
}

func (s *VendorService) CreateVendor(vendor *models.Vendor) error {
	return s.db.Create(vendor).Error
}

func (s *VendorService) GetVendorByName(name string) (*models.Vendor, error) {
	var vendor models.Vendor
	if err := s.db.Where("name = ?", name).First(&vendor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("vendor %q not found", name)
		}
		return nil, err
	}

	return &vendor, nil
}

func (s *VendorService) SelectBestVendor(vendors []models.Vendor) (*models.VendorSelectionResponse, error) {
	if len(vendors) == 0 {
		return nil, errors.New("at least one vendor is required")
	}

	minPrice := vendors[0].Price
	maxPrice := vendors[0].Price
	const MaxTrust = 5.0
	for _, vendor := range vendors {
		if vendor.Price <= 0 {
			return nil, fmt.Errorf("vendor %q has invalid price", vendor.Name)
		}
		if vendor.Trust < 0 || vendor.Trust > 5 {
			return nil, fmt.Errorf("vendor %q trust must be between 0 and 5", vendor.Name)
		}
		if vendor.Price < minPrice {
			minPrice = vendor.Price
		}
		if vendor.Price > maxPrice {
			maxPrice = vendor.Price
		}
	}

	scored := make([]models.ScoredVendor, 0, len(vendors))
	for _, vendor := range vendors {
		priceScore := normalizePriceScore(vendor.Price, minPrice, maxPrice)
		trustScore := float64(vendor.Trust) / MaxTrust
		finalScore := (trustScore * 0.65) + (priceScore * 0.35)

		scored = append(scored, models.ScoredVendor{
			Vendor: vendor,
			Score:  math.Round(finalScore*1000) / 1000,
		})
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	best := scored[0]
	reason := fmt.Sprintf(
		"%s selected for its strong reliability (%.1f/5 trust score) and competitive pricing (%.2f), weighted at 65%% trust and 35%% price.",
		best.Vendor.Name,
		best.Vendor.Trust,
		best.Vendor.Price,
	)

	return &models.VendorSelectionResponse{
		RecommendedVendor: best.Vendor,
		Score:             best.Score,
		Reason:            reason,
		ScoredVendors:     scored,
	}, nil
}

func normalizePriceScore(price, minPrice, maxPrice float64) float64 {
	if minPrice == maxPrice {
		return 1
	}

	return (maxPrice - price) / (maxPrice - minPrice)
}
