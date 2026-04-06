package models

import "time"

type Order struct {
	ID          string    `json:"id" gorm:"primaryKey;size:32"`
	Vendor      string    `json:"vendor" gorm:"not null"`
	VendorID    uint      `json:"vendorId" gorm:"index;not null"`
	VendorRef   Vendor    `json:"-" gorm:"foreignKey:VendorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Amount      float64   `json:"amount" gorm:"not null"`
	Status      string    `json:"status" gorm:"not null"`
	PaymentTxID string    `json:"paymentTxId,omitempty"`
	QR          *QR       `json:"qr,omitempty" gorm:"foreignKey:OrderID;references:ID"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreateOrderRequest struct {
	Vendor string  `json:"vendor" binding:"required"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type PaymentActionRequest struct {
	OrderID string `json:"orderId" binding:"required"`
}

type ConfirmDeliveryRequest struct {
	OrderID string `json:"orderId" binding:"required"`
}

type PaymentActionResponse struct {
	OrderID string `json:"orderId"`
	TxID    string `json:"txID"`
	Status  string `json:"status"`
}

type ConfirmDeliveryResponse struct {
	Message string `json:"message"`
	TxID    string `json:"txID"`
	Order   *Order `json:"order"`
	Status  string `json:"status"`
}
