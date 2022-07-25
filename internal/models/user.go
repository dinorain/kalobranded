package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	UserRoleAdmin = "admin"
	UserRoleUser  = "user"
)

// User model
type User struct {
	UserID          uuid.UUID `json:"user_id" db:"user_id"`
	Email           string    `json:"email" db:"email"`
	FirstName       string    `json:"first_name" db:"first_name"`
	LastName        string    `json:"last_name" db:"last_name"`
	DeliveryAddress string    `json:"delivery_address" db:"delivery_address"`
	Role            string    `json:"role" db:"role"`
	Avatar          *string   `json:"avatar" db:"avatar"`
	Password        string    `json:"-" db:"password"`
	CreatedAt       time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

func (u *User) SanitizePassword() {
	u.Password = ""
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) ComparePasswords(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) PrepareCreate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Password = strings.TrimSpace(u.Password)
	u.DeliveryAddress = strings.TrimSpace(u.DeliveryAddress)

	if err := u.HashPassword(); err != nil {
		return err
	}

	if u.Role != "" {
		u.Role = strings.ToLower(strings.TrimSpace(u.Role))
		if u.Role != UserRoleAdmin && u.Role != UserRoleUser {
			return fmt.Errorf("role invalid: %v", u.Role)
		}
	}

	return nil
}

// Get avatar string
func (u *User) GetAvatar() string {
	if u.Avatar == nil {
		return ""
	}
	return *u.Avatar
}
