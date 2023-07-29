package main

import (
	"context"
	"time"

	"github.com/jailtonjunior94/go-rabbitmq/pkg/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	OrdersExchange = "orders"
)

var (
	Exchanges = []*rabbitmq.Exchange{
		rabbitmq.NewExchange(OrdersExchange, "direct"),
	}

	Bindings = []*rabbitmq.Binding{
		rabbitmq.NewBinding("order", OrdersExchange),
		rabbitmq.NewBinding("logs", OrdersExchange),
		rabbitmq.NewBinding("finance_orders", OrdersExchange),
	}
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

	_, err = rabbitmq.NewAmqpBuilder(channel).
		DeclareExchanges(Exchanges...).
		DeclareBindings(Bindings...).
		DeclarePrefetchCount(5).
		WithDlq().
		WithRetry().
		DeclareTtl(30 * time.Second).
		Apply()

	if err != nil {
		panic(err)
	}

	err = channel.PublishWithContext(context.Background(), OrdersExchange, "", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("abc"),
	})

	if err != nil {
		panic(err)
	}
}
