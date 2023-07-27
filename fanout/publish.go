package main

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Fanout - envia cópia de mensagens para várias filas
func main() {
	ctx := context.Background()
	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer channel.Close()

	orderQueue, err := channel.QueueDeclare("order", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	logsQueue, err := channel.QueueDeclare("logs", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	financeQueue, err := channel.QueueDeclare("finance_orders", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	err = channel.ExchangeDeclare("order", "fanout", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	_ = channel.QueueBind(orderQueue.Name, "", "order", false, nil)
	_ = channel.QueueBind(logsQueue.Name, "", "order", false, nil)
	_ = channel.QueueBind(financeQueue.Name, "", "order", false, nil)

	err = channel.PublishWithContext(ctx, "order", "", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("Hello World"),
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully Published Message to Queue")
}
