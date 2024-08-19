package connection

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// NotifyBlockedSafe safely wraps the (*amqp.Connection).NotifyBlocked method
func (m *Manager) startNotifyBlockedHandler() {
	m.connectionMu.RLock()
	defer m.connectionMu.RUnlock()
	receiver := m.connection.NotifyBlocked(make(chan amqp.Blocking))
	for r := range receiver {
		m.blockedMu.Lock()
		if r.Active {
			m.logger.Warnf("server TCP blocking")
			m.blocked = true
		} else {
			m.blocked = false
			m.logger.Warnf("server TCP unblocking")
		}
		m.blockedMu.Unlock()
	}
	return
}

func (m *Manager) Blocked() bool {
	m.blockedMu.RLock()
	defer m.blockedMu.RUnlock()
	return m.blocked
}
