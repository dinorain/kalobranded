package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/dinorain/checkoutaja/internal/models"
)

type SellerResponseDto struct {
	SellerID      uuid.UUID `json:"seller_id"`
	Email         string    `json:"email"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Avatar        *string   `json:"avatar"`
	PickupAddress string    `json:"pickup_address"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func SellerResponseFromModel(seller *models.Seller) *SellerResponseDto {
	return &SellerResponseDto{
		SellerID:      seller.SellerID,
		Email:         seller.Email,
		FirstName:     seller.FirstName,
		LastName:      seller.LastName,
		Avatar:        seller.Avatar,
		PickupAddress: seller.PickupAddress,
		CreatedAt:     seller.CreatedAt,
		UpdatedAt:     seller.UpdatedAt,
	}
}
