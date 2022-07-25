package dto

type BrandUpdateRequestDto struct {
	BrandName     *string `json:"brand_name" validate:"omitempty,lte=30"`
	Logo          *string `json:"logo" validate:"omitempty"`
	PickupAddress *string `json:"pickup_address" validate:"omitempty"`
}
