package dto

import "github.com/dinorain/kalobranded/pkg/utils"

type ProductFindResponseDto struct {
	Meta utils.PaginationMetaDto `json:"meta"`
	Data interface{}             `json:"data"`
}
