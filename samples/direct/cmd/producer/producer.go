package main

import (
	"context"
	"encoding/json"
	"math/rand"

	"time"

	"github.com/jailtonjunior94/go-rabbitmq/pkg/rabbitmq"
	"github.com/jailtonjunior94/go-rabbitmq/samples/direct/internal/entities"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	OrdersExchange = "order"
	OrderCreated   = "order_created"
	OrderUpdated   = "order_updated"
	OrderQueue     = "order"
	FinanceQueue   = "finance_order"
)

var (
	Exchanges = []*rabbitmq.Exchange{
		rabbitmq.NewExchange(OrdersExchange, "direct"),
	}

	Bindings = []*rabbitmq.Binding{
		rabbitmq.NewBindingRouting(OrderQueue, OrdersExchange, OrderCreated),
		rabbitmq.NewBindingRouting(OrderQueue, OrdersExchange, OrderUpdated),
		rabbitmq.NewBindingRouting(FinanceQueue, OrdersExchange, OrderCreated),
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
		DeclareTtl(60 * time.Second).
		Apply()

	if err != nil {
		panic(err)
	}

	order := entities.NewOrder(rand.Int63(), rand.Float64())
	payload, _ := json.Marshal(order)

	if err := channel.PublishWithContext(context.Background(), OrdersExchange, OrderCreated, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        payload,
	}); err != nil {
		panic(err)
	}

	order.Update(rand.Float64())
	payloadUpdated, _ := json.Marshal(order)

	if err := channel.PublishWithContext(context.Background(), OrdersExchange, OrderUpdated, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        payloadUpdated,
	}); err != nil {
		panic(err)
	}
}
