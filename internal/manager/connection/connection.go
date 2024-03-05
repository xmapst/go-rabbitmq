package connection

import (
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/xmapst/go-rabbitmq/internal/dispatcher"
	"github.com/xmapst/go-rabbitmq/internal/logger"
)

// Manager -
type Manager struct {
	logger               logger.Logger
	url                  string
	connection           *amqp.Connection
	amqpConfig           amqp.Config
	connectionMux        *sync.RWMutex
	ReconnectInterval    time.Duration
	reconnectionCount    uint
	reconnectionCountMux *sync.Mutex
	dispatcher           *dispatcher.Dispatcher
}

// New creates a new connection manager
func New(url string, conf amqp.Config, log logger.Logger, reconnectInterval time.Duration) (*Manager, error) {
	conn, err := amqp.DialConfig(url, conf)
	if err != nil {
		return nil, err
	}
	connManager := Manager{
		logger:               log,
		url:                  url,
		connection:           conn,
		amqpConfig:           conf,
		connectionMux:        &sync.RWMutex{},
		ReconnectInterval:    reconnectInterval,
		reconnectionCount:    0,
		reconnectionCountMux: &sync.Mutex{},
		dispatcher:           dispatcher.New(),
	}
	go connManager.startNotifyClose()
	return &connManager, nil
}

// Close safely closes the current channel and connection
func (m *Manager) Close() error {
	m.logger.Infof("closing connection manager...")
	m.connectionMux.Lock()
	defer m.connectionMux.Unlock()

	err := m.connection.Close()
	if err != nil {
		m.logger.Errorf("close err: %v", err)
		return err
	}
	m.logger.Infof("amqp connection closed gracefully")
	return nil
}

// NotifyReconnect adds a new subscriber that will receive error messages whenever
// the connection manager has successfully reconnected to the server
func (m *Manager) NotifyReconnect() (<-chan error, chan<- struct{}) {
	return m.dispatcher.AddSubscriber()
}

// CheckoutConnection -
func (m *Manager) CheckoutConnection() *amqp.Connection {
	m.connectionMux.RLock()
	return m.connection
}

// CheckinConnection -
func (m *Manager) CheckinConnection() {
	m.connectionMux.RUnlock()
}

// startNotifyCancelOrClosed listens on the channel's cancelled and closed
// notifiers. When it detects a problem, it attempts to reconnect.
// Once reconnected, it sends an error back on the manager's notifyCancelOrClose
// channel
func (m *Manager) startNotifyClose() {
	notifyCloseChan := m.connection.NotifyClose(make(chan *amqp.Error, 1))

	err := <-notifyCloseChan
	if err != nil {
		m.logger.Errorf("attempting to reconnect to amqp server after connection close with error: %v", err)
		m.reconnectLoop()
		m.logger.Warnf("successfully reconnected to amqp server")
		if _err := m.dispatcher.Dispatch(err); _err != nil {
			m.logger.Warnf("connection dispatch err: %v", err)
		}
	}
	if err == nil {
		m.logger.Infof("amqp connection closed gracefully")
	}
}

// GetReconnectionCount -
func (m *Manager) GetReconnectionCount() uint {
	m.reconnectionCountMux.Lock()
	defer m.reconnectionCountMux.Unlock()
	return m.reconnectionCount
}

func (m *Manager) incrementReconnectionCount() {
	m.reconnectionCountMux.Lock()
	defer m.reconnectionCountMux.Unlock()
	m.reconnectionCount++
}

// reconnectLoop continuously attempts to reconnect
func (m *Manager) reconnectLoop() {
	for {
		m.logger.Infof("waiting %s seconds to attempt to reconnect to amqp server", m.ReconnectInterval)
		time.Sleep(m.ReconnectInterval)
		err := m.reconnect()
		if err != nil {
			m.logger.Errorf("error reconnecting to amqp server: %v", err)
		} else {
			m.incrementReconnectionCount()
			go m.startNotifyClose()
			return
		}
	}
}

// reconnect safely closes the current channel and obtains a new one
func (m *Manager) reconnect() error {
	m.connectionMux.Lock()
	defer m.connectionMux.Unlock()
	newConn, err := amqp.DialConfig(m.url, m.amqpConfig)
	if err != nil {
		return err
	}

	if err = m.connection.Close(); err != nil {
		m.logger.Warnf("error closing connection while reconnecting: %v", err)
	}

	m.connection = newConn
	return nil
}
