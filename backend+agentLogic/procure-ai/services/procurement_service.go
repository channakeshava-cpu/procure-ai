package services

import "procure-ai/models"

type ProcurementService struct {
	orderService      *OrderService
	blockchainService *BlockchainService
}

func NewProcurementService(orderService *OrderService, blockchainService *BlockchainService) *ProcurementService {
	return &ProcurementService{
		orderService:      orderService,
		blockchainService: blockchainService,
	}
}

func (s *ProcurementService) CreateOrder(req models.CreateOrderRequest) (*models.Order, error) {
	return s.orderService.CreateOrder(req)
}

func (s *ProcurementService) LockFunds(orderID string) (*models.PaymentActionResponse, error) {
	if _, err := s.orderService.GetOrder(orderID); err != nil {
		return nil, err
	}

	txID, err := s.blockchainService.LockFunds(orderID)
	if err != nil {
		return nil, err
	}

	order, err := s.orderService.MarkFundsLocked(orderID)
	if err != nil {
		return nil, err
	}

	return &models.PaymentActionResponse{
		OrderID: order.ID,
		TxID:    txID,
		Status:  order.Status,
	}, nil
}

func (s *ProcurementService) ReleasePayment(orderID string) (*models.PaymentActionResponse, error) {
	if _, err := s.orderService.GetOrder(orderID); err != nil {
		return nil, err
	}

	txID, err := s.blockchainService.ReleasePayment(orderID)
	if err != nil {
		return nil, err
	}

	order, err := s.orderService.MarkPaymentReleased(orderID)
	if err != nil {
		return nil, err
	}

	return &models.PaymentActionResponse{
		OrderID: order.ID,
		TxID:    txID,
		Status:  order.Status,
	}, nil
}

func (s *ProcurementService) ConfirmDelivery(orderID string) (*models.ConfirmDeliveryResponse, error) {
	order, err := s.orderService.ConfirmDelivery(orderID)
	if err != nil {
		return nil, err
	}

	txID, err := s.blockchainService.ReleasePayment(orderID)
	if err != nil {
		return nil, err
	}

	order, err = s.orderService.MarkPaymentReleased(orderID)
	if err != nil {
		return nil, err
	}

	return &models.ConfirmDeliveryResponse{
		Message: "delivery confirmed and payment released",
		TxID:    txID,
		Order:   order,
		Status:  order.Status,
	}, nil
}
