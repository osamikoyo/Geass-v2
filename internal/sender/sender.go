package sender

import (
	"github.com/osamikoyo/geass-v2/pkg/loger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Sender struct {
	Channel *amqp.Channel
	Queue amqp.Queue
	Logger loger.Logger
}