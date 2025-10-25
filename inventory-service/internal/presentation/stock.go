package presentation

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type StockRequest struct {
	ProductCode string `json:"product_code"`
	Amount      int32  `json:"amount"`
}

func (r StockRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ProductCode, validation.Required),
		validation.Field(&r.Amount, validation.Required, validation.Min(1)),
	)
}

type CreateStockRequest struct {
	Stocks []StockRequest `json:"stocks"`
}

type StockResponse struct {
	IsAvailable bool `json:"is_available"`
}
