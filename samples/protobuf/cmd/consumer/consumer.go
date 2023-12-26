package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jailtonjunior94/go-rabbitmq/pkg/rabbitmq"
	"github.com/jailtonjunior94/go-rabbitmq/samples/protobuf/pkg/order"

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

func handlerMessage(ctx context.Context, body []byte) error {
	order := &order.OrderMessage{}
	if err := proto.Unmarshal(body, order); err != nil {
		log.Fatalln("Failed to parse order:", err)
	}

	return nil
}

func main() {
	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}

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

	ctx := context.Background()

	consumer, err := rabbitmq.NewConsumer(
		rabbitmq.WithConnection(connection),
		rabbitmq.WithChannel(channel),
		rabbitmq.WithQueue(FinanceQueue),
		rabbitmq.WithHandler(handlerMessage),
	)

	consumer.Consume(ctx)

	forever := make(chan bool)

	fmt.Println("Successfully connected to our RabbitMQ instance")
	fmt.Println(" [*] - Waiting for messages")
	<-forever
}
