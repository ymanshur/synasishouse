package presentation

import (
	"time"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateProductRequest struct {
	Code  string `json:"code"`
	Total int32  `json:"total"`
}

func (r CreateProductRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Code, validation.Required),
		validation.Field(&r.Total, validation.Required, validation.Min(1)),
	)
}

type GetProductRequest struct {
	ID string `json:"id"`
}

func (r GetProductRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID))
}

type UpdateProductRequest struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}

func (r UpdateProductRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID),
		validation.Field(&r.Code, validation.Required),
	)
}

type ProductResponse struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Total     int32     `json:"total"`
	Hold      int32     `json:"hold"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
