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

func (h *OrderHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var req presentation.OrderRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Warn().Err(err).Msg("cannot bind request")
		response.New().
			WithCode(http.StatusBadRequest).
			WithMessage(err.Error()).
			JSON(c)
		return
	}

	res, err := h.orderUseCase.Create(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("cannot create order")
		response.New().
			WithTranslationError(err).
			JSON(c)
		return
	}

	if !res.IsAvailable {
		response.New().
			WithCode(http.StatusUnprocessableEntity).
			WithMessage("stock is unavailable").
			JSON(c)
		return
	}

	response.New().
		WithCode(http.StatusOK).
		WithMessage("stock is available").
		JSON(c)
}
