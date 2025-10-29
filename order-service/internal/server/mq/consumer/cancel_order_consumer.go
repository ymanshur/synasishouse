package consumer

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
	"github.com/ymanshur/synasishouse/order/internal/config"
	"github.com/ymanshur/synasishouse/order/internal/presentation"
	"github.com/ymanshur/synasishouse/order/internal/typex"
	"github.com/ymanshur/synasishouse/order/internal/usecase"
)

type CancelOrderConsumer struct {
	config       config.RabbitMQClientConfig
	orderUseCase usecase.Orderer
}

func NewCancelOrder(config config.RabbitMQClientConfig, orderUseCase usecase.Orderer) *CancelOrderConsumer {
	return &CancelOrderConsumer{
		config:       config,
		orderUseCase: orderUseCase,
	}
}

func (c *CancelOrderConsumer) Config() config.RabbitMQClientConfig {
	return c.config
}

func (c *CancelOrderConsumer) Process(ctx context.Context, msg *amqp091.Delivery) {
	var req presentation.UpdateOrderRequest
	err := json.Unmarshal(msg.Body, &req)
	if err != nil {
		log.Warn().Err(err).Msg("cannot unmarshal request")
		msg.Ack(false)
		return
	}

	_, err = c.orderUseCase.Cancel(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("cannot cancel order")

		var unavailableErr typex.Unavailable
		if errors.As(err, &unavailableErr) {
			msg.Nack(false, true)
			return
		}

		msg.Ack(false)
		return
	}
}
