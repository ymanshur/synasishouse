package bootstrap

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
	"github.com/ymanshur/synasishouse/order/internal/appctx"
)

// RegistryRabbitMQ establish RabbitMQ connection
func RegistryRabbitMQ(config appctx.RabbitMQConfig) *amqp091.Connection {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.VHost,
	)

	conn, err := amqp091.Dial(url)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to rabbitmq")
	}

	return conn
}
