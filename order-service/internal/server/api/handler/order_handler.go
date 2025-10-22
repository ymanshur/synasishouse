package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/ymanshur/synasishouse/order/internal/presentation"
	"github.com/ymanshur/synasishouse/order/internal/server/api/response"
	"github.com/ymanshur/synasishouse/order/internal/usecase"
)

type OrderHandler struct {
	orderUseCase usecase.Orderer
}

func NewOrder(orderUseCase usecase.Orderer) *OrderHandler {
	return &OrderHandler{orderUseCase: orderUseCase}
}

func (h *OrderHandler) Checkout(c *gin.Context) {
	ctx := c.Request.Context()

	var req presentation.OrderRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Warn().Err(err)
		response.New().
			WithCode(http.StatusBadRequest).
			WithErrors(err).
			JSON(c)
		return
	}

	ok, err := h.orderUseCase.Checkout(ctx, req)
	if err != nil {
		log.Err(err)
		response.New().
			WithTranslationError(err).
			JSON(c)
		return
	}

	if !ok {
		log.Warn().Err(err)
		response.New().
			WithCode(http.StatusUnprocessableEntity).
			WithMessage("stock is unavailable").
			JSON(c)
	}

	response.New().
		WithCode(http.StatusOK).
		WithMessage("stock is available").
		JSON(c)
}
