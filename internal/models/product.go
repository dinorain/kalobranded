package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Product model
type Product struct {
	ProductID   uuid.UUID `json:"product_id" db:"product_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Price       float64   `json:"price" db:"price"`
	BrandID    uuid.UUID `json:"brand_id" db:"brand_id"`
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

func (p *Product) PrepareCreate() error {
	p.Name = strings.TrimSpace(p.Name)
	p.Description = strings.TrimSpace(p.Description)
	return nil
}
