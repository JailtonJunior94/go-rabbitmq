package main

import (
	"context"
	"fmt"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
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

	messages, err := channel.Consume("logs", "", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)
	go func() {
		for message := range messages {
			err := processMessage(message.Body)
			if err != nil {
				handleRetry(channel, message)
				continue
			}
			message.Ack(true)
		}
	}()

	fmt.Println("Successfully connected to our RabbitMQ instance")
	fmt.Println(" [*] - Waiting for messages")
	<-forever
}

func processMessage(message []byte) error {
	_, err := strconv.Atoi(string(message))
	if err != nil {
		return err
	}
	return nil
}

func handleRetry(ch *amqp.Channel, delivery amqp.Delivery) {
	if retry(delivery) {
		delivery.Nack(false, false)
		return
	}

	sendDLQ(ch, delivery)
	delivery.Ack(true)
}

func retry(delivery amqp.Delivery) bool {
	deaths, ok := delivery.Headers["x-death"].([]interface{})
	if !ok || len(deaths) <= 0 {
		return true
	}

	for _, death := range deaths {
		values, ok := death.(amqp.Table)
		if !ok {
			return true
		}
		if count, ok := values["count"].(int64); !ok || count < 3 {
			return true
		}
	}

	return false
}

func sendDLQ(ch *amqp.Channel, delivery amqp.Delivery) {
	delete(delivery.Headers, "x-death")
	err := ch.PublishWithContext(context.Background(), "", "logs.dlq", false, false, amqp.Publishing{
		ContentType: delivery.ContentType,
		Body:        delivery.Body,
	})
	if err != nil {
		panic(err)
	}
}
