package Consumer

import (
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"listener-srv/cmd/handler"
	"listener-srv/internal/entity/eventEntity"
	"listener-srv/internal/eventService/event"
	"log"
)

type Consumer interface {
	Listen(topics []string) error
	HandlePayload(payload eventEntity.Payload)
}

type amqpConsumer struct {
	conn      *amqp091.Connection
	queueName string
	handler   *handler.MailHandler
}

func (a amqpConsumer) HandlePayload(payload eventEntity.Payload) {
	switch payload.Action {
	case "email":
		if payload.SubAction == "" {
			log.Println("send batch email")
		}
		log.Println("send email here")
		// use mail-grpc to call the mailing service
		data, err := a.handler.SendMail(payload)
		if err != nil {
			log.Println(data.Result)
			return
		}
	}
}

func (a amqpConsumer) Listen(topics []string) error {
	ch, err := a.conn.Channel()
	if err != nil {
		log.Println(err)
		return err
	}

	queue, err := event.DeclareRandomQueue(ch)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, topic := range topics {
		err = ch.QueueBind(
			queue.Name,
			topic,
			event.ExchangeName,
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for msg := range messages {
			var payload eventEntity.Payload
			json.Unmarshal(msg.Body, &payload)
			//log.Printf("#%v", payload)
			go a.HandlePayload(payload)
		}
	}()

	log.Printf("\n waiting on message on [Exchange, Queue] [mail_exchange, %s]\n", queue.Name)
	<-forever
	return nil
}

func NewAmqpConsumer(conn *amqp091.Connection, handler *handler.MailHandler) (Consumer, error) {

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	err = event.DeclareExchange(channel)
	if err != nil {
		return nil, err
	}

	return &amqpConsumer{
		conn:    conn,
		handler: handler,
	}, err
}
