package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ymanshur/synasishouse/pkg/util"
)

// Response presentation contract object
type Response struct {
	Code    int `json:"code"`
	Errors  any `json:"errors,omitempty"`
	Data    any `json:"data,omitempty"`
	Meta    any `json:"meta,omitempty"`
	Message any `json:"message,omitempty"`
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
	return &Response{}
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
	r.Code = code
	return r
}

// WithData set response data
func (r *Response) WithData(v any) *Response {
	r.Data = v
	return r
}

// WithTranslationError translates response error.
// It will assign the status code, Message and/or Errors.
func (r *Response) WithTranslationError(err error) *Response {
	switch {
	case errors.As(err, &validationErrs):
		r.Code = http.StatusUnprocessableEntity
		r.Errors = convertValidationErrors(validationErrs)
	case errors.As(err, &unprocessableEntityErr):
		r.Code = http.StatusUnprocessableEntity
		r.Message = unprocessableEntityErr.Error()
	case errors.As(err, &conflictErr):
		r.Code = http.StatusConflict
		r.Message = conflictErr.Error()
	case errors.As(err, &notFoundErr):
		r.Code = http.StatusNotFound
		r.Message = notFoundErr.Error()
	default:
		if code, ok := util.TranslateGRPCError(err); ok {
			r.Code = code
			r.Message = err.Error()
			return r
		}

		r.Code = http.StatusInternalServerError
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
	c.JSON(r.Code, r)
}
