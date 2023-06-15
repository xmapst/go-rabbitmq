package connection

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// NotifyBlockedSafe safely wraps the (*amqp.Connection).NotifyBlocked method
func (m *Manager) NotifyBlockedSafe(
	receiver chan amqp.Blocking,
) chan amqp.Blocking {
	m.connectionMux.RLock()
	defer m.connectionMux.RUnlock()

	return m.connection.NotifyBlocked(
		receiver,
	)
}
