package connection

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// NotifyBlockedSafe safely wraps the (*amqp.Connection).NotifyBlocked method
func (m *Manager) NotifyBlockedSafe(
	receiver chan amqp.Blocking,
) chan amqp.Blocking {
	m.connectionMu.RLock()
	defer m.connectionMu.RUnlock()
	return m.connection.NotifyBlocked(
		receiver,
	)
}
