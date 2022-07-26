package dto

import (
	"github.com/google/uuid"
)

type BrandRegisterRequestDto struct {
	BrandName     string  `json:"brand_name" validate:"required,lte=30"`
	PickupAddress string  `json:"pickup_address" validate:"required"`
	Logo          *string `json:"logo"`
}

type BrandRegisterResponseDto struct {
	BrandID uuid.UUID `json:"user_id" validate:"required"`
}
