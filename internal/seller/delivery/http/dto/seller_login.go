package dto

import (
	"github.com/google/uuid"
)

type SellerLoginRequestDto struct {
	Email    string `json:"email" validate:"required,lte=60,email"`
	Password string `json:"password" validate:"required"`
}

type SellerLoginResponseDto struct {
	SellerID uuid.UUID                      `json:"user_id" validate:"required"`
	Tokens   *SellerRefreshTokenResponseDto `json:"tokens" validate:"required"`
}
