package connection

import (
	"errors"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/xmapst/go-rabbitmq/internal/dispatcher"
	"github.com/xmapst/go-rabbitmq/internal/logger"
)

// Manager -
type Manager struct {
	logger              logger.Logger
	resolver            Resolver
	connection          *amqp.Connection
	amqpConfig          amqp.Config
	connectionMu        *sync.RWMutex
	ReconnectInterval   time.Duration
	reconnectionCount   uint
	reconnectionCountMu *sync.Mutex
	dispatcher          *dispatcher.Dispatcher

	blocked   bool
	blockedMu *sync.RWMutex
}

type Resolver interface {
	Resolve() ([]string, error)
}

// dial will attempt to connect to the a list of urls in the order they are
// given.
func dial(log logger.Logger, resolver Resolver, conf amqp.Config) (*amqp.Connection, error) {
	urls, err := resolver.Resolve()
	if err != nil {
		return nil, fmt.Errorf("error resolving amqp server urls: %w", err)
	}

	var errs []error
	for _, url := range urls {
		conn, err := amqp.DialConfig(url, amqp.Config(conf))
		if err == nil {
			return conn, err
		}
		log.Warnf("failed to connect to amqp server %s: %v", url, err)
		errs = append(errs, err)
	}
	return nil, errors.Join(errs...)
}

// New creates a new connection manager
func New(resolver Resolver, conf amqp.Config, log logger.Logger, reconnectInterval time.Duration) (*Manager, error) {
	conn, err := dial(log, resolver, conf)
	if err != nil {
		return nil, err
	}
	connManager := Manager{
		logger:              log,
		resolver:            resolver,
		connection:          conn,
		amqpConfig:          conf,
		connectionMu:        &sync.RWMutex{},
		ReconnectInterval:   reconnectInterval,
		reconnectionCount:   0,
		reconnectionCountMu: &sync.Mutex{},
		dispatcher:          dispatcher.New(),
		blocked:             false,
		blockedMu:           &sync.RWMutex{},
	}
	go connManager.startNotifyClose()
	go connManager.startNotifyBlockedHandler()
	return &connManager, nil
}

// Close safely closes the current channel and connection
func (m *Manager) Close() error {
	m.logger.Infof("closing connection manager...")
	m.connectionMu.Lock()
	defer m.connectionMu.Unlock()

	err := m.connection.Close()
	if err != nil {
		return err
	}
	return nil
}

// NotifyReconnect adds a new subscriber that will receive error messages whenever
// the connection manager has successfully reconnected to the server
func (m *Manager) NotifyReconnect() (<-chan error, chan<- struct{}) {
	return m.dispatcher.AddSubscriber()
}

// CheckoutConnection -
func (m *Manager) CheckoutConnection() *amqp.Connection {
	m.connectionMu.RLock()
	return m.connection
}

// CheckinConnection -
func (m *Manager) CheckinConnection() {
	m.connectionMu.RUnlock()
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
		_ = m.dispatcher.Dispatch(err)
	}
	if err == nil {
		m.logger.Infof("amqp connection closed gracefully")
	}
}

// GetReconnectionCount -
func (m *Manager) GetReconnectionCount() uint {
	m.reconnectionCountMu.Lock()
	defer m.reconnectionCountMu.Unlock()
	return m.reconnectionCount
}

func (m *Manager) incrementReconnectionCount() {
	m.reconnectionCountMu.Lock()
	defer m.reconnectionCountMu.Unlock()
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
			go m.startNotifyBlockedHandler()
			return
		}
	}
}

// reconnect safely closes the current channel and obtains a new one
func (m *Manager) reconnect() error {
	m.connectionMu.Lock()
	defer m.connectionMu.Unlock()

	conn, err := dial(m.logger, m.resolver, m.amqpConfig)
	if err != nil {
		return err
	}

	if err = m.connection.Close(); err != nil {
		m.logger.Warnf("error closing connection while reconnecting: %v", err)
	}

	m.connection = conn
	return nil
}
