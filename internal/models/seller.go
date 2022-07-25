package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Brand model
type Brand struct {
	BrandID       uuid.UUID `json:"brand_id" db:"brand_id"`
	BrandName     string    `json:"brand_name" db:"brand_name"`
	PickupAddress string    `json:"pickup_address" db:"pickup_address"`
	Logo          *string   `json:"logo" db:"logo"`
	CreatedAt     time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

func (s *Brand) PrepareCreate() error {
	s.BrandName = strings.TrimSpace(s.BrandName)
	s.PickupAddress = strings.TrimSpace(s.PickupAddress)

	return nil
}
