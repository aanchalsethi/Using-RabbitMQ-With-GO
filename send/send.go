package main

import (
	"context"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, key string) {
	if err != nil {
		log.Fatalf("%s: %s", key, err)
	}
}

func main() {
	// connect to rabbitmq
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "failed to connect to rabbitmq")

	//close the connection
	defer conn.Close()

	// declare a channel
	ch, err := conn.Channel()
	failOnError(err, "failed to open a channel")

	//close the channel
	defer ch.Close()

	//declare the queue
	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World!"
	err = ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)

	failOnError(err, "Failed to publish a message")
	log.Printf("[x] Sent %s\n", body)
}
