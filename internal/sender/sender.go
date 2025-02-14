package sender

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/osamikoyo/geass-v2/internal/models"
	"github.com/osamikoyo/geass-v2/pkg/config"
	"github.com/osamikoyo/geass-v2/pkg/loger"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
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

func (s *Sender) Send(content models.PageInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := sonic.Marshal(&content)
	if err != nil{
		return fmt.Errorf("cant marshal body: %w", err)
	}

	err = s.Channel.PublishWithContext(
		ctx,
		"",
		s.Queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body: body,
			},
		)
	if err != nil{
		return fmt.Errorf("cant publishing: %w", err)
	}

	return nil
}