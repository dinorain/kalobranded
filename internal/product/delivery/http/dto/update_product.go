package dto

type ProductUpdateRequestDto struct {
	Name        *string  `json:"name" validate:"omitempty,lte=30"`
	Description *string  `json:"description" validate:"omitempty,lte=250"`
	Price       *float64 `json:"price" validate:"omitempty"`
}
