package mq

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
	"github.com/ymanshur/synasishouse/order/internal/appctx"
	"github.com/ymanshur/synasishouse/order/internal/server"
	"github.com/ymanshur/synasishouse/order/internal/server/mq/consumer"
	"github.com/ymanshur/synasishouse/order/internal/usecase"
	"golang.org/x/sync/errgroup"
)

type Processor interface {
	Process(ctx context.Context, msg *amqp091.Delivery)
}

type Consumer interface {
	Config() appctx.RabbitMQClientConfig
	Processor
}

type mqServer struct {
	conn      *amqp091.Connection
	consumers []Consumer
}

func NewServer(conn *amqp091.Connection, orderUseCase usecase.Orderer) server.Server {
	config := appctx.LoadConfig()
	configConsumer := config.RabbitMQClient

	consumers := []Consumer{
		consumer.NewCancelOrder(configConsumer.CancelOrder, orderUseCase),
	}

	return &mqServer{
		conn:      conn,
		consumers: consumers,
	}
}

func (s *mqServer) Run(ctx context.Context) error {
	wg, wgCtx := errgroup.WithContext(ctx)

	for _, consumer := range s.consumers {
		c := consumer
		wg.Go(func() (err error) {
			config := c.Config()

			channel, err := s.conn.Channel()
			if err != nil {
				return fmt.Errorf("open server channel: %w", err)
			}

			_, err = channel.QueueDeclare(config.Queue, true, false, false, false, nil)
			if err != nil {
				return fmt.Errorf("declare queue: %w", err)
			}

			err = channel.QueueBind(config.Queue, config.Key, config.Exchange, false, nil)
			if err != nil {
				return fmt.Errorf("bind queue: %w", err)
			}

			// TODO: Consume RabbitMQ queue by multiple workers

			return consume(wgCtx, channel, config.Queue, c)
		})
	}
	return wg.Wait()
}

func consume(ctx context.Context, ch *amqp091.Channel, queue string, processor Processor) error {
	msgs, err := ch.ConsumeWithContext(ctx, queue, fmt.Sprintf("consumer-%s", queue), false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume queue message: %w", err)
	}

	for msg := range msgs {
		process(ctx, processor, &msg)
	}

	return nil
}

func process(ctx context.Context, processor Processor, msg *amqp091.Delivery) {
	defer func() {
		if r := recover(); r != nil {
			log.Error().
				Any("panic", r).
				Bytes("stack", debug.Stack()).
				Msg("recovered from panic")

			msg.Ack(false)
		}
	}()

	processor.Process(ctx, msg)
}
