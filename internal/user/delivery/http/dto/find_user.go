package dto

import "github.com/dinorain/checkoutaja/pkg/utils"

type UserFindResponseDto struct {
	Meta utils.PaginationMetaDto `json:"meta"`
	Data interface{}             `json:"data"`
}
