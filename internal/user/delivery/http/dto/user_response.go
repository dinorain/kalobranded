package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/dinorain/checkoutaja/internal/models"
)

type UserResponseDto struct {
	UserID          uuid.UUID `json:"user_id"`
	Email           string    `json:"email"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Role            string    `json:"role"`
	Avatar          *string   `json:"avatar"`
	DeliveryAddress string    `json:"delivery_address"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func UserResponseFromModel(user *models.User) *UserResponseDto {
	return &UserResponseDto{
		UserID:          user.UserID,
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Role:            user.Role,
		Avatar:          user.Avatar,
		DeliveryAddress: user.DeliveryAddress,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}
}
