package channel

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ConsumeSafe safely wraps the (*amqp.Channel).Consume method
func (m *Manager) ConsumeSafe(
	queue,
	consumer string,
	autoAck,
	exclusive,
	noLocal,
	noWait bool,
	args amqp.Table,
) (<-chan amqp.Delivery, error) {
	m.channelMux.RLock()
	defer m.channelMux.RUnlock()

	return m.channel.Consume(
		queue,
		consumer,
		autoAck,
		exclusive,
		noLocal,
		noWait,
		args,
	)
}

// QueueDeclarePassiveSafe safely wraps the (*amqp.Channel).QueueDeclarePassive method
func (m *Manager) QueueDeclarePassiveSafe(
	name string,
	durable bool,
	autoDelete bool,
	exclusive bool,
	noWait bool,
	args amqp.Table,
) (amqp.Queue, error) {
	m.channelMux.RLock()
	defer m.channelMux.RUnlock()

	return m.channel.QueueDeclarePassive(
		name,
		durable,
		autoDelete,
		exclusive,
		noWait,
		args,
	)
}

// QueueDeclareSafe safely wraps the (*amqp.Channel).QueueDeclare method
func (m *Manager) QueueDeclareSafe(
	name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table,
) (amqp.Queue, error) {
	m.channelMux.RLock()
	defer m.channelMux.RUnlock()

	return m.channel.QueueDeclare(
		name,
		durable,
		autoDelete,
		exclusive,
		noWait,
		args,
	)
}

// ExchangeDeclarePassiveSafe safely wraps the (*amqp.Channel).ExchangeDeclarePassive method
func (m *Manager) ExchangeDeclarePassiveSafe(
	name string, kind string, durable bool, autoDelete bool, internal bool, noWait bool, args amqp.Table,
) error {
	m.channelMux.RLock()
	defer m.channelMux.RUnlock()

	return m.channel.ExchangeDeclarePassive(
		name,
		kind,
		durable,
		autoDelete,
		internal,
		noWait,
		args,
	)
}

// ExchangeDeclareSafe safely wraps the (*amqp.Channel).ExchangeDeclare method
func (m *Manager) ExchangeDeclareSafe(
	name string, kind string, durable bool, autoDelete bool, internal bool, noWait bool, args amqp.Table,
) error {
	m.channelMux.RLock()
	defer m.channelMux.RUnlock()

	return m.channel.ExchangeDeclare(
		name,
		kind,
		durable,
		autoDelete,
		internal,
		noWait,
		args,
	)
}

func (m *Manager) ExchangeBindSafe(
	destination, key, source string, noWait bool, args amqp.Table,
) error {
	m.channelMux.RLock()
	defer m.channelMux.RUnlock()

	return m.channel.ExchangeBind(
		destination,
		key,
		source,
		noWait,
		args,
	)
}

// QueueBindSafe safely wraps the (*amqp.Channel).QueueBind method
func (m *Manager) QueueBindSafe(
	name string, key string, exchange string, noWait bool, args amqp.Table,
) error {
	m.channelMux.RLock()
	defer m.channelMux.RUnlock()

	return m.channel.QueueBind(
		name,
		key,
		exchange,
		noWait,
		args,
	)
}

// QosSafe safely wraps the (*amqp.Channel).Qos method
func (m *Manager) QosSafe(
	prefetchCount int, prefetchSize int, global bool,
) error {
	m.channelMux.RLock()
	defer m.channelMux.RUnlock()

	return m.channel.Qos(
		prefetchCount,
		prefetchSize,
		global,
	)
}

/*
PublishSafe safely wraps the (*amqp.Channel).Publish method.
*/
func (m *Manager) PublishSafe(
	exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing,
) error {
	m.channelMux.RLock()
	defer m.channelMux.RUnlock()

	return m.channel.PublishWithContext(
		context.Background(),
		exchange,
		key,
		mandatory,
		immediate,
		msg,
	)
}

// PublishWithContextSafe safely wraps the (*amqp.Channel).PublishWithContext method.
func (m *Manager) PublishWithContextSafe(
	ctx context.Context, exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing,
) error {
	m.channelMux.RLock()
	defer m.channelMux.RUnlock()

	return m.channel.PublishWithContext(
		ctx,
		exchange,
		key,
		mandatory,
		immediate,
		msg,
	)
}

func (m *Manager) PublishWithDeferredConfirmWithContextSafe(
	ctx context.Context, exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing,
) (*amqp.DeferredConfirmation, error) {
	m.channelMux.RLock()
	defer m.channelMux.RUnlock()

	return m.channel.PublishWithDeferredConfirmWithContext(
		ctx,
		exchange,
		key,
		mandatory,
		immediate,
		msg,
	)
}

// NotifyReturnSafe safely wraps the (*amqp.Channel).NotifyReturn method
func (m *Manager) NotifyReturnSafe(
	c chan amqp.Return,
) chan amqp.Return {
	m.channelMux.RLock()
	defer m.channelMux.RUnlock()

	return m.channel.NotifyReturn(
		c,
	)
}

// ConfirmSafe safely wraps the (*amqp.Channel).Confirm method
func (m *Manager) ConfirmSafe(
	noWait bool,
) error {
	m.channelMux.Lock()
	defer m.channelMux.Unlock()

	return m.channel.Confirm(
		noWait,
	)
}

// NotifyPublishSafe safely wraps the (*amqp.Channel).NotifyPublish method
func (m *Manager) NotifyPublishSafe(
	confirm chan amqp.Confirmation,
) chan amqp.Confirmation {
	m.channelMux.RLock()
	defer m.channelMux.RUnlock()

	return m.channel.NotifyPublish(
		confirm,
	)
}

// NotifyFlowSafe safely wraps the (*amqp.Channel).NotifyFlow method
func (m *Manager) NotifyFlowSafe(
	c chan bool,
) chan bool {
	m.channelMux.RLock()
	defer m.channelMux.RUnlock()

	return m.channel.NotifyFlow(
		c,
	)
}
