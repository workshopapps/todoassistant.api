package Emitter

import (
	"context"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"test-va/internals/entity/eventEntity"
	"test-va/internals/msg-queue/event"
	"time"
)

type Emitter interface {
	Push(payload eventEntity.Payload, severity string) error
}

type amqpEmitter struct {
	conn *amqp091.Connection
}

func (a amqpEmitter) Push(payload eventEntity.Payload, severity string) error {
	channel, err := a.conn.Channel()
	if err != nil {
		return err
	}
	//defer channel.Close()

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	log.Println("Pushing to Chanel")
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*30)
	defer cancelFunc()

	err = channel.PublishWithContext(ctx, event.ExchangeName, severity, false, false, amqp091.Publishing{
		ContentType: "text/plain",
		Body:        data,
	})
	if err != nil {
		return err
	}
	return nil
}

func NewAmqpEmitter(conn *amqp091.Connection) (Emitter, error) {

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer channel.Close()

	err = event.DeclareExchange(channel)
	if err != nil {
		return nil, err
	}

	return &amqpEmitter{conn: conn}, nil
}
