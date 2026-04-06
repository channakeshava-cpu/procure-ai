package services

import (
	"encoding/base64"
	"fmt"

	qrcode "github.com/skip2/go-qrcode"
	"procure-ai/models"

	"gorm.io/gorm"
)

type QRService struct {
	db           *gorm.DB
	orderService *OrderService
}

func NewQRService(db *gorm.DB, orderService *OrderService) *QRService {
	return &QRService{db: db, orderService: orderService}
}

func (s *QRService) GenerateQR(orderID string) (*models.GenerateQRResponse, error) {
	if _, err := s.orderService.GetOrder(orderID); err != nil {
		return nil, err
	}

	qrContent := fmt.Sprintf("PROCURE-ORDER:%s", orderID)
	png, err := qrcode.Encode(qrContent, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}

	record := &models.QR{
		OrderID: orderID,
		QRCode:  qrContent,
		QRImage: base64.StdEncoding.EncodeToString(png),
	}
	if err := s.db.Where(models.QR{OrderID: orderID}).Assign(models.QR{
		QRCode:  record.QRCode,
		QRImage: record.QRImage,
	}).FirstOrCreate(record).Error; err != nil {
		return nil, err
	}

	return &models.GenerateQRResponse{
		OrderID: orderID,
		QRCode:  qrContent,
		QRImage: record.QRImage,
		Message: "QR generated successfully",
	}, nil
}

func (s *QRService) VerifyQR(orderID, qrCode string) (*models.VerifyQRResponse, error) {
	if _, err := s.orderService.GetOrder(orderID); err != nil {
		return nil, err
	}

	var record models.QR
	if err := s.db.First(&record, "order_id = ?", orderID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &models.VerifyQRResponse{
				OrderID: orderID,
				Valid:   false,
				Message: "QR not found for order",
			}, nil
		}
		return nil, err
	}

	expected := record.QRCode
	valid := qrCode == expected
	message := "QR verification failed"
	if valid {
		message = "QR verified successfully"
	}

	return &models.VerifyQRResponse{
		OrderID: orderID,
		Valid:   valid,
		Message: message,
	}, nil
}
