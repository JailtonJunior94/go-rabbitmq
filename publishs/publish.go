package main

import (
	"fmt"

	"github.com/streadway/amqp"
)

func main() {
	fmt.Println("Server is running - Golang RabbitMQ")

	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer connection.Close()

	fmt.Println("Successfully connected to our RabbitMQ Instance")

	channel, err := connection.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare("golang-queue", true, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(queue)

	err = channel.Publish("", "golang-queue", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("Hello World"),
	})

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("Successfully Published Message to Queue")
}
