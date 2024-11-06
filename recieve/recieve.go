package main

import (
	"log"

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

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "failed to register a consumer")

	var forever chan struct{}
	// iterate over msgs
	go func() {
		for d := range msgs {
			log.Printf("recieved a message: %s", d.Body)
		}
	}()

	log.Printf("[*] Waiting for messages. To exit press CTRL +C ")
	<-forever
}
