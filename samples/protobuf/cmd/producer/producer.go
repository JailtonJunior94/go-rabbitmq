package main

import (
	"context"
	"log"
	"math/rand"

	"time"

	"github.com/jailtonjunior94/go-rabbitmq/pkg/rabbitmq"
	"github.com/jailtonjunior94/go-rabbitmq/samples/protobuf/pkg/order"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
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

	order := &order.OrderMessage{Id: uuid.New().String(), Amount: rand.Float32()}
	body, err := proto.Marshal(order)
	if err != nil {
		log.Fatalln("Failed to encode address book:", err)
	}

	if err := channel.PublishWithContext(context.Background(), OrdersExchange, OrderCreated, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	}); err != nil {
		panic(err)
	}
}
