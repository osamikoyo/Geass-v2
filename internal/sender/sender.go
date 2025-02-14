package sender

import (
	"fmt"
	"github.com/osamikoyo/geass-v2/pkg/config"
	"github.com/osamikoyo/geass-v2/pkg/loger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Sender struct {
	Channel *amqp.Channel
	Queue amqp.Queue
	Logger loger.Logger
}

func New(cfg *config.Config) (*Sender, error) {
	conn, err := amqp.Dial(cfg.AmqpConnectUrl)
	if err != nil{
		return nil, fmt.Errorf("cant connect to amqp: %w",err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil{
		return nil, fmt.Errorf("cant get channel: %w", err)
	}
	defer ch.Close()

	que, err := ch.QueueDeclare(
		"content",
			false,
			false,
			false,
			false,
			nil,
		)

	return &Sender{
		Channel: ch,
		Queue: que,
	}, nil
}