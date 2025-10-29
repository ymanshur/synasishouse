package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

var _ Producer = (*AutoCancelOrderProducer)(nil)

type AutoCancelOrderProducer struct {
	channel  *amqp091.Channel
	exchange string
	key      string
}

func NewAutoCancelOrderProducer(
	conn *amqp091.Connection,
	exchange string,
	key string,
) (*AutoCancelOrderProducer, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("open channel: %w", err)
	}

	args := amqp091.Table{
		"x-delayed-type": "direct",
	}
	if err := channel.ExchangeDeclare(exchange, "x-delayed-message", true, false, false, false, args); err != nil {
		if amqpErr, ok := err.(*amqp091.Error); ok && amqpErr.Code != amqp091.PreconditionFailed {
			return nil, fmt.Errorf("declare exchange %s: %w", exchange, err)
		}
	}

	return &AutoCancelOrderProducer{
		channel:  channel,
		exchange: exchange,
		key:      key,
	}, nil
}

func (p *AutoCancelOrderProducer) Send(ctx context.Context, msg any) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	return p.channel.PublishWithContext(
		ctx,
		p.exchange,
		p.key,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Headers: amqp091.Table{
				"x-delay": (24 * time.Hour).Milliseconds(),
			},
			Body: body,
		},
	)
}
