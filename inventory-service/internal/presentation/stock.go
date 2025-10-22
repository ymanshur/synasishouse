package presentation

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type StockRequest struct {
	Code   string `json:"code"`
	Amount int32  `json:"amount"`
}

func (r StockRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Code, validation.Required),
		validation.Field(&r.Amount, validation.Required, validation.Min(1)),
	)
}

type StockResponse struct {
	IsAvailable bool `json:"is_available"`
}
