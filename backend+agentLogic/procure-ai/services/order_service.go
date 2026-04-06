package services

import (
	"errors"
	"fmt"

	"procure-ai/models"

	"gorm.io/gorm"
)

type OrderService struct {
	db            *gorm.DB
	vendorService *VendorService
}

func NewOrderService(db *gorm.DB, vendorService *VendorService) *OrderService {
	return &OrderService{
		db:            db,
		vendorService: vendorService,
	}
}

func (s *OrderService) CreateOrder(req models.CreateOrderRequest) (*models.Order, error) {
	if req.Vendor == "" {
		return nil, errors.New("vendor is required")
	}
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	vendor, err := s.vendorService.GetVendorByName(req.Vendor)
	if err != nil {
		return nil, err
	}

	var count int64
	if err := s.db.Model(&models.Order{}).Count(&count).Error; err != nil {
		return nil, err
	}

	orderID := fmt.Sprintf("ORD-%04d", count+1)
	order := &models.Order{
		ID:       orderID,
		Vendor:   vendor.Name,
		VendorID: vendor.ID,
		Amount:   req.Amount,
		Status:   "created",
	}

	if err := s.db.Create(order).Error; err != nil {
		return nil, err
	}

	return s.GetOrder(orderID)
}

func (s *OrderService) GetOrder(orderID string) (*models.Order, error) {
	var order models.Order
	if err := s.db.Preload("QR").First(&order, "id = ?", orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order %s not found", orderID)
		}
		return nil, err
	}

	return &order, nil
}

func (s *OrderService) ConfirmDelivery(orderID string) (*models.Order, error) {
	if err := s.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", "delivered").Error; err != nil {
		return nil, err
	}

	return s.GetOrder(orderID)
}

func (s *OrderService) MarkFundsLocked(orderID string) (*models.Order, error) {
	if err := s.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", "funds_locked").Error; err != nil {
		return nil, err
	}

	return s.GetOrder(orderID)
}

func (s *OrderService) MarkPaymentReleased(orderID string) (*models.Order, error) {
	if err := s.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", "payment_released").Error; err != nil {
		return nil, err
	}

	return s.GetOrder(orderID)
}
