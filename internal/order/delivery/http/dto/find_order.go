package dto

import "github.com/dinorain/checkoutaja/pkg/utils"

type OrderFindResponseDto struct {
	Meta utils.PaginationMetaDto `json:"meta"`
	Data interface{}             `json:"data"`
}
