package rabbitmq

import (
	"errors"

	"github.com/xmapst/go-rabbitmq/internal/manager/channel"
)

type Declarator struct {
	chanManager *channel.Manager
}

func NewDeclarator(conn *Conn) (*Declarator, error) {
	if conn.connManager == nil {
		return nil, errors.New("connection manager can't be nil")
	}

	chanManager, err := channel.New(conn.connManager, &stdDebugLogger{}, conn.connManager.ReconnectInterval)
	if err != nil {
		return nil, err
	}

	result := &Declarator{
		chanManager: chanManager,
	}

	return result, nil
}

func (d *Declarator) Close() {
	_ = d.chanManager.Close()
}

func (d *Declarator) Exchange(optionFuncs ...func(*PublisherOptions)) error {
	defaultOptions := getDefaultPublisherOptions()
	options := &defaultOptions
	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}

	return declareExchange(d.chanManager, options.ExchangeOptions)
}

func (d *Declarator) Queue(queue string, optionFuncs ...func(*ConsumerOptions)) error {
	defaultOptions := getDefaultConsumerOptions(queue)
	options := &defaultOptions
	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}

	return declareQueue(d.chanManager, options.QueueOptions)
}

func (d *Declarator) BindExchanges(bindings []ExchangeBinding) error {
	for _, binding := range bindings {
		err := d.chanManager.ExchangeBindSafe(
			binding.Destination,
			binding.RoutingKey,
			binding.Source,
			binding.NoWait,
			tableToAMQPTable(binding.Args),
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Declarator) BindQueues(bindings []QueueBinding) error {
	for _, binding := range bindings {
		err := d.chanManager.QueueBindSafe(
			binding.Queue,
			binding.RoutingKey,
			binding.Exchange,
			binding.NoWait,
			tableToAMQPTable(binding.Args),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func declareQueue(chanManager *channel.Manager, options QueueOptions) error {
	if !options.Declare {
		return nil
	}
	if options.Passive {
		_, err := chanManager.QueueDeclarePassiveSafe(
			options.Name,
			options.Durable,
			options.AutoDelete,
			options.Exclusive,
			options.NoWait,
			tableToAMQPTable(options.Args),
		)
		if err != nil {
			return err
		}
		return nil
	}
	_, err := chanManager.QueueDeclareSafe(
		options.Name,
		options.Durable,
		options.AutoDelete,
		options.Exclusive,
		options.NoWait,
		tableToAMQPTable(options.Args),
	)
	if err != nil {
		return err
	}
	return nil
}

func declareExchange(chanManager *channel.Manager, options ExchangeOptions) error {
	if !options.Declare {
		return nil
	}
	if options.Passive {
		err := chanManager.ExchangeDeclarePassiveSafe(
			options.Name,
			options.Kind,
			options.Durable,
			options.AutoDelete,
			options.Internal,
			options.NoWait,
			tableToAMQPTable(options.Args),
		)
		if err != nil {
			return err
		}
		return nil
	}
	err := chanManager.ExchangeDeclareSafe(
		options.Name,
		options.Kind,
		options.Durable,
		options.AutoDelete,
		options.Internal,
		options.NoWait,
		tableToAMQPTable(options.Args),
	)
	if err != nil {
		return err
	}
	return nil
}

func declareSuperExchange(chanManager *channel.Manager, options SuperExchangeOptions) error {
	if !options.Declare {
		return nil
	}
	if options.Passive {
		err := chanManager.ExchangeDeclarePassiveSafe(
			options.Name,
			options.Kind,
			options.Durable,
			options.AutoDelete,
			options.Internal,
			options.NoWait,
			tableToAMQPTable(options.Args),
		)
		if err != nil {
			return err
		}
		return nil
	}
	err := chanManager.ExchangeDeclareSafe(
		options.Name,
		options.Kind,
		options.Durable,
		options.AutoDelete,
		options.Internal,
		options.NoWait,
		tableToAMQPTable(options.Args),
	)
	if err != nil {
		return err
	}
	return nil
}

func declareBindings(chanManager *channel.Manager, options ConsumerOptions) error {
	for _, binding := range options.Bindings {
		if !binding.Declare {
			continue
		}
		err := chanManager.QueueBindSafe(
			options.QueueOptions.Name,
			binding.RoutingKey,
			options.ExchangeOptions.Name,
			binding.NoWait,
			tableToAMQPTable(binding.Args),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func declareSuperBindings(chanManager *channel.Manager, options SuperConsumerOptions) error {
	for _, exchangeOptions := range options.ExchangeOptionsSlice {
		for _, binding := range exchangeOptions.Bindings {
			if !binding.Declare {
				continue
			}
			err := chanManager.QueueBindSafe(
				options.QueueOptions.Name,
				binding.RoutingKey,
				exchangeOptions.Name,
				binding.NoWait,
				tableToAMQPTable(binding.Args),
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
