package dto

import (
	db "github.com/ymanshur/synasishouse/inventory/db/sqlc"
	"github.com/ymanshur/synasishouse/inventory/internal/presentation"
)

func ProductToResponse(i db.Product) presentation.ProductResponse {
	return presentation.ProductResponse{
		ID:        i.ID.String(),
		Code:      i.Code,
		Total:     i.Total,
		Reserved:  i.Reserved,
		UpdatedAt: i.UpdatedAt,
		CreatedAt: i.CreatedAt,
	}
}
