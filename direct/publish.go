package main

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Direct - envia mensagens para filas que estiver com routing key iguais
func main() {
	ctx := context.Background()
	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
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

	financeQueue, err := channel.QueueDeclare("finance_order", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	err = channel.ExchangeDeclare("order", "direct", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	_ = channel.QueueBind(orderQueue.Name, "order_new", "order", false, nil)
	_ = channel.QueueBind(orderQueue.Name, "order_upd", "order", false, nil)
	_ = channel.QueueBind(financeQueue.Name, "order_new", "order", false, nil)

	if err := channel.PublishWithContext(ctx, "order", "order_new", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("Hello World"),
	}); err != nil {
		panic(err)
	}

	if err := channel.PublishWithContext(ctx, "order", "order_upd", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("Hello World"),
	}); err != nil {
		panic(err)
	}
}
