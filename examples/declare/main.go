package main

import (
	"log"

	"github.com/xmapst/go-rabbitmq"
)

func main() {
	conn, err := rabbitmq.NewConn(
		"amqp://guest:guest@localhost",
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	declare, err := rabbitmq.NewDeclarator(conn)
	if err != nil {
		log.Fatal(err)
	}
	defer declare.Close()

	// declare dead-letter queue
	err = declare.Queue(
		"re_my_queue",
		rabbitmq.WithConsumerOptionsRoutingKey("my_routing_key"),
		rabbitmq.WithConsumerOptionsExchangeName("events"),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		log.Fatal(err)
	}
	err = declare.BindQueues([]rabbitmq.QueueBinding{
		{
			Exchange:   "events",
			RoutingKey: "my_routing_key",
			Queue:      "re_my_queue",
			NoWait:     false,
			Args: map[string]interface{}{
				"x-dead-letter-exchange":    "events",
				"x-dead-letter-routing-key": "my_routing_key",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
