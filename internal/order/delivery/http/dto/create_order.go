package dto

import (
	"github.com/google/uuid"
)

type OrderCreateRequestDto struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  uint64    `json:"quantity" validate:"required"`
}

type OrderCreateResponseDto struct {
	OrderID uuid.UUID `json:"order_id" validate:"required"`
}
