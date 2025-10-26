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

	response.New().
		WithCode(http.StatusOK).
		WithMessage("order created successfuly").
		WithData(res).
		JSON(c)
}

func (h *OrderHandler) Settle(c *gin.Context) {
	ctx := c.Request.Context()

	var req presentation.UpdateOrderRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Warn().Err(err).Msg("cannot bind request")
		response.New().
			WithCode(http.StatusBadRequest).
			WithMessage(err.Error()).
			JSON(c)
		return
	}

	req.OrderNo = c.Param("order_no")

	res, err := h.orderUseCase.Settle(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("cannot settle order")
		response.New().
			WithTranslationError(err).
			JSON(c)
		return
	}

	response.New().
		WithCode(http.StatusOK).
		WithMessage("order settled successfuly").
		WithData(res).
		JSON(c)
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	ctx := c.Request.Context()

	var req presentation.UpdateOrderRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Warn().Err(err).Msg("cannot bind request")
		response.New().
			WithCode(http.StatusBadRequest).
			WithMessage(err.Error()).
			JSON(c)
		return
	}

	req.OrderNo = c.Param("order_no")

	res, err := h.orderUseCase.Cancel(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("cannot cancel order")
		response.New().
			WithTranslationError(err).
			JSON(c)
		return
	}

	response.New().
		WithCode(http.StatusOK).
		WithMessage("order canceled successfuly").
		WithData(res).
		JSON(c)
}
