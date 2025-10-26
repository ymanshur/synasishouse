package presentation

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type OrderRequest struct {
	OrderNo string               `json:"order_no"`
	UserID  string               `json:"user_id"`
	Details []OrderDetailRequest `json:"details"`
}

func (r OrderRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.OrderNo, validation.Required),
		validation.Field(&r.UserID, validation.Required, is.UUID),
		validation.Field(&r.Details),
	)
}

type OrderDetailRequest struct {
	ProductCode string `json:"product_code"`
	Amount      int32  `json:"amount"`
}

func (r OrderDetailRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ProductCode, validation.Required),
		validation.Field(&r.Amount, validation.Required, validation.Min(1)),
	)
}

type OrderResponse struct {
	OrderNo string                `json:"order_no"`
	UserID  string                `json:"user_id"`
	Status  string                `json:"status"`
	Details []OrderDetailResponse `json:"details"`
}

type OrderDetailResponse struct {
	ProductCode string `json:"product_code"`
	Amount      int32  `json:"amount"`
}

type UpdateOrderRequest struct {
	OrderNo string
	UserID  string `json:"user_id"`
}

func (r UpdateOrderRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.OrderNo, validation.Required),
		validation.Field(&r.UserID, validation.Required, is.UUID),
	)
}

type UpdateOrderResponse struct {
	OrderNo string `json:"order_no"`
	UserID  string `json:"user_id"`
	Status  string `json:"status"`
}
