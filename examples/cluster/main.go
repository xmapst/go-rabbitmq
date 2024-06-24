package main

import (
	"log"

	"github.com/xmapst/go-rabbitmq"
)

func main() {
	resolver := rabbitmq.NewStaticResolver(
		[]string{
			"amqp://guest:guest@host1",
			"amqp://guest:guest@host2",
			"amqp://guest:guest@host3",
		},
		false, /* shuffle */
	)

	conn, err := rabbitmq.NewClusterConn(resolver)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

}
