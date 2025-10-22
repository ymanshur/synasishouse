package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ymanshur/synasishouse/pkg/util"
)

// Response presentation contract object
type Response struct {
	code    int
	Success bool `json:"success"`
	Errors  any  `json:"errors,omitempty"`
	Data    any  `json:"data,omitempty"`
	Meta    any  `json:"meta,omitempty"`
	Message any  `json:"message,omitempty"`
}

// MetaData represent meta data response for list data
type MetaData struct {
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	TotalPage  int64 `json:"total_page"`
	TotalCount int64 `json:"total_count"`
}

// New return [Response] instance
func New() *Response {
	return &Response{Success: true}
}

// NewMeta return [MetaData] instance
func NewMeta(page, limit, totalPage, totalCount int64) MetaData {
	return MetaData{
		Page:       page,
		Limit:      limit,
		TotalPage:  totalPage,
		TotalCount: totalCount,
	}
}

// WithCode set response status code
func (r *Response) WithCode(code int) *Response {
	r.code = code
	return r
}

// GetCode set response status code
func (r *Response) GetCode() int {
	return r.code
}

// WithData set response data
func (r *Response) WithData(v any) *Response {
	r.Data = v
	return r
}

// WithErrors parse the Gin validation errors
func (r *Response) WithErrors(err error) *Response {
	r.Success = false

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errs := []validationError{}
		for _, f := range validationErrors {
			errs = append(errs,
				validationError{
					Field: f.Field(),
					Msg:   validationMsgFromFieldError(f),
				})
		}
		r.Errors = errs
		r.Message = "invalid paramater"
		return r
	}
	r.Message = err.Error()
	return r
}

// WithTranslationError translates response error.
// It will assign the status code, message and/or error.
func (r *Response) WithTranslationError(err error) *Response {
	r.Success = false

	switch {
	case errors.As(err, &unprocessableEntityErr):
		r.code = http.StatusUnprocessableEntity
		r.Message = unprocessableEntityErr.Error()
	case errors.As(err, &conflictErr):
		r.code = http.StatusConflict
		r.Message = conflictErr.Error()
	case errors.As(err, &notFoundErr):
		r.code = http.StatusNotFound
		r.Message = notFoundErr.Error()
	default:
		if code, ok := util.TranslateGRPCError(err); ok {
			r.code = code
			r.Message = err.Error()
			return r
		}

		r.code = http.StatusInternalServerError
		r.Message = "something went wrong"
	}

	return r
}

// WithMeta set response meta data
func (r *Response) WithMeta(v MetaData) *Response {
	r.Meta = v
	return r
}

// WithMessage set response message
func (r *Response) WithMessage(msg string) *Response {
	if msg != "" {
		r.Message = msg
	}

	return r
}

// JSON render response through [*gin.Context]
func (r *Response) JSON(c *gin.Context) {
	c.JSON(r.code, r)
}
