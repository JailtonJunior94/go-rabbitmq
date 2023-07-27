package main

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	ctx := context.Background()
	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	setup(channel)
	producer(ctx, channel)
	consumer(channel)
}

func setup(channel *amqp.Channel) {
	err := channel.ExchangeDeclare("dlq_exchange", "fanout", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	dqlQueue, err := channel.QueueDeclare("dlq_queue", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	_ = channel.QueueBind(dqlQueue.Name, "", "dlq_exchange", false, nil)

	exchangeArgs := make(amqp.Table)
	exchangeArgs["x-dead-letter-exchange"] = "dlq_exchange"
	_, err = channel.QueueDeclare("task_queue", true, false, false, false, exchangeArgs)
	if err != nil {
		panic(err)
	}
}

func producer(ctx context.Context, channel *amqp.Channel) {
	if err := channel.PublishWithContext(ctx, "task_queue", "", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("Hello World"),
	}); err != nil {
		panic(err)
	}
}

func consumer(channel *amqp.Channel) {
	messages, err := channel.Consume("task_queue", "", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)
	go func() {
		for message := range messages {
			fmt.Printf("Recieved Message: %s\n", message.Body)
			message.Ack(false)
		}
	}()

	fmt.Println("Successfully connected to our RabbitMQ instance")
	fmt.Println(" [*] - Waiting for messages")
	<-forever
}
