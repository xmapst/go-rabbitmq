package rabbitmq

type ExchangeBinding struct {
	Destination string
	RoutingKey  string
	Source      string
	Args        Table
	NoWait      bool
}

type QueueBinding struct {
	Exchange   string
	RoutingKey string
	Queue      string
	Args       Table
	NoWait     bool
}
