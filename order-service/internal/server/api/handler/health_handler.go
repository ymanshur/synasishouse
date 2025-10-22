package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ymanshur/synasishouse/order/internal/server/api/response"
)

type HealthHandler struct{}

func NewHealth() *HealthHandler {
	return &HealthHandler{}
}

// Check the service healthy
func (u *HealthHandler) Check(c *gin.Context) {
	rsp := response.New().
		WithCode(http.StatusOK).
		WithMessage("UP")
	c.JSON(rsp.GetCode(), rsp)
}
