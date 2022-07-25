package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/dinorain/checkoutaja/internal/models"
)

type ProductResponseDto struct {
	ProductID   uuid.UUID `json:"product_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	SellerID    uuid.UUID `json:"seller_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ProductResponseFromModel(user *models.Product) *ProductResponseDto {
	return &ProductResponseDto{
		ProductID:   user.ProductID,
		Name:        user.Name,
		Description: user.Description,
		Price:       user.Price,
		SellerID:    user.SellerID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
