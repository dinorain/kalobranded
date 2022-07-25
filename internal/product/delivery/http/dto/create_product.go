package dto

import (
	"github.com/google/uuid"
)

type ProductCreateRequestDto struct {
	Name        string  `json:"name" validate:"required,lte=30"`
	Description string  `json:"description" validate:"required,lte=250"`
	Price       float64 `json:"price" validate:"required"`
}

type ProductCreateResponseDto struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
}
