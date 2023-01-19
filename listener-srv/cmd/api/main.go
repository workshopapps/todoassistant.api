package main

import (
	"listener-srv/cmd/handler"
	"listener-srv/internal/eventService/Consumer"
	"log"
	"math"
	"os"
	"time"

	ampq "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect to rabbit mq
	conn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	mailHandler, err := handler.NewMailHandler()
	if err != nil {
		panic(err)
	}
	defer mailHandler.CloseConn()
	// create a consumer
	consumer, err := Consumer.NewAmqpConsumer(conn, mailHandler)
	if err != nil {
		log.Println(err)
		log.Panicf("error consumer cannot be nil")
	}

	// watch the queue and consume the events
	err = consumer.Listen([]string{"info", "email", "log"})
	if err != nil {
		log.Println("error")
		return
	}
}

func connect() (*ampq.Connection, error) {
	var count int64
	var backoff = 1 * time.Second
	var connection *ampq.Connection

	// don't continue until connects

	for {
		c, err := ampq.Dial("amqp://guest:guest@localhost:5672")
		if err != nil {
			log.Println("Not Connected yet....")
			count++
		} else {
			connection = c
			break
		}

		if count > 5 {
			log.Println(" error connecting", nil)
			return nil, err
		} else {
			backoff = time.Duration(math.Pow(float64(count), 2)) * time.Second
			log.Println("waiting....")
			time.Sleep(backoff)
			continue
		}
	}
	log.Println("connected to rabbitMQ")
	return connection, nil
}
