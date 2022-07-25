package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/dinorain/checkoutaja/internal/models"
)

type OrderResponseDto struct {
	OrderID                    uuid.UUID        `json:"order_id"`
	UserID                     uuid.UUID        `json:"user_id"`
	SellerID                   uuid.UUID        `json:"seller_id"`
	Item                       models.OrderItem `json:"item"`
	Quantity                   uint64           `json:"quantity"`
	TotalPrice                 float64          `json:"total_price"`
	Status                     string           `json:"status"`
	DeliverySourceAddress      string           `json:"delivery_source_address"`
	DeliveryDestinationAddress string           `json:"delivery_destination_address"`
	CreatedAt                  time.Time        `json:"created_at,omitempty"`
	UpdatedAt                  time.Time        `json:"updated_at,omitempty"`
}

func OrderResponseFromModel(order *models.Order) *OrderResponseDto {
	return &OrderResponseDto{
		OrderID:                    order.OrderID,
		UserID:                     order.UserID,
		SellerID:                   order.SellerID,
		Item:                       order.Item,
		Quantity:                   order.Quantity,
		TotalPrice:                 order.TotalPrice,
		Status:                     order.Status,
		DeliverySourceAddress:      order.DeliverySourceAddress,
		DeliveryDestinationAddress: order.DeliveryDestinationAddress,
		CreatedAt:                  order.CreatedAt,
		UpdatedAt:                  order.UpdatedAt,
	}
}
