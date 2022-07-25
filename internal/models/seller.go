package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Seller model
type Seller struct {
	SellerID      uuid.UUID `json:"seller_id" db:"seller_id"`
	Email         string    `json:"email" db:"email"`
	FirstName     string    `json:"first_name" db:"first_name"`
	LastName      string    `json:"last_name" db:"last_name"`
	PickupAddress string    `json:"pickup_address" db:"pickup_address"`
	Avatar        *string   `json:"avatar" db:"avatar"`
	Password      string    `json:"-" db:"password"`
	CreatedAt     time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

func (s *Seller) SanitizePassword() {
	s.Password = ""
}

func (s *Seller) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	s.Password = string(hashedPassword)
	return nil
}

func (s *Seller) ComparePasswords(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(s.Password), []byte(password))
}

func (s *Seller) PrepareCreate() error {
	s.Email = strings.ToLower(strings.TrimSpace(s.Email))
	s.Password = strings.TrimSpace(s.Password)
	s.PickupAddress = strings.TrimSpace(s.PickupAddress)

	if err := s.HashPassword(); err != nil {
		return err
	}
	return nil
}

// Get avatar string
func (s *Seller) GetAvatar() string {
	if s.Avatar == nil {
		return ""
	}
	return *s.Avatar
}
