package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/dinorain/kalobranded/internal/models"
)

type BrandResponseDto struct {
	BrandID       uuid.UUID `json:"brand_id"`
	BrandName     string    `json:"brand_name"`
	Logo          *string   `json:"logo"`
	PickupAddress string    `json:"pickup_address"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func BrandResponseFromModel(brand *models.Brand) *BrandResponseDto {
	return &BrandResponseDto{
		BrandID:       brand.BrandID,
		BrandName:     brand.BrandName,
		Logo:          brand.Logo,
		PickupAddress: brand.PickupAddress,
		CreatedAt:     brand.CreatedAt,
		UpdatedAt:     brand.UpdatedAt,
	}
}
