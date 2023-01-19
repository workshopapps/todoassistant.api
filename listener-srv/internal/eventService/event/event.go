package event

import "github.com/rabbitmq/amqp091-go"

const ExchangeName = "mail_exchange"

func DeclareExchange(ch *amqp091.Channel) error {
	return ch.ExchangeDeclare(
		ExchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
}

func DeclareRandomQueue(ch *amqp091.Channel) (amqp091.Queue, error) {
	return ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
}
