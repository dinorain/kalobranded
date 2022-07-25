package dto

type SellerRefreshTokenDto struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type SellerRefreshTokenResponseDto struct {
	AccessToken  string `json:"access_token" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}
