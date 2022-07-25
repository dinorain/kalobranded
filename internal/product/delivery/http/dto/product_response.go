package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/dinorain/kalobranded/internal/models"
)

type ProductResponseDto struct {
	ProductID   uuid.UUID `json:"product_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	BrandID    uuid.UUID `json:"brand_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ProductResponseFromModel(user *models.Product) *ProductResponseDto {
	return &ProductResponseDto{
		ProductID:   user.ProductID,
		Name:        user.Name,
		Description: user.Description,
		Price:       user.Price,
		BrandID:    user.BrandID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
