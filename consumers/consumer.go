package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	fmt.Println("Consumer Application")

	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer channel.Close()

	messages, err := channel.Consume("golang-queue", "", true, false, false, false, nil)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			fmt.Printf("Recieved Message: %s\n", d.Body)
		}
	}()

	fmt.Println("Successfully connected to our RabbitMQ instance")
	fmt.Println(" [*] - Waiting for messages")
	<-forever
}
