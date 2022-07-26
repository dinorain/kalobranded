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
	BrandID     uuid.UUID `json:"brand_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ProductResponseFromModel(product *models.Product) *ProductResponseDto {
	return &ProductResponseDto{
		ProductID:   product.ProductID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		BrandID:     product.BrandID,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}
