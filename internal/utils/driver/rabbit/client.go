// TODO: Try to simplify.
package rabbit

import (
	"strings"

	"github.com/assembla/cony"
)

// PublishFunc describes the publish to RabbitMQ function.
type PublishFunc func(*cony.Publisher)

// Client is a wrapper around cony.Client which keeps track of the RabbitMQ
// connections.
type Client struct {
	config *Config
	cony   *cony.Client
}

// NewDB creates new connection to RabbitMQ using cony.
func NewClient(cfg *Config) (*Client, error) {
	if err := createEnv(&cfg.Server); err != nil {
		return nil, err
	}

	client := cony.NewClient(cony.URL(cfg.Server.Server), cony.Backoff(cony.DefaultBackoff))

	return &Client{config: cfg, cony: client}, nil
}

// Dialer returns a pointer to the Dialer with which the connection was made.
func (client *Client) Config() *Config {
	return client.config
}

// Cony returns pointer to cony.Client.
func (client *Client) Cony() *cony.Client {
	return client.cony
}

// Errors returns errors channel.
// If you do not read the channel, then when its buffer is full, the client will
// stop working.
func (client *Client) Errors() <-chan error {
	return client.cony.Errors()
}

// Loop returns true on successful receipt of data from the queue.
// Loop should be run as condition for `for` with receiving from (*Client).Errors().
//
// It will manage AMQP connection, run queue and exchange declarations, consumers.
// Will start to return false once (*Client).Close() called.
func (client *Client) Loop() bool {
	return client.cony.Loop()
}

// NewPublisherProcess creates new *cony.Publisher to exchange by routing_key
// and starts publishing process described in cb.
// NewPublisherProcess does not declare the queue, so if it was not created earlier,
// the publication will go to the void.
func (client *Client) NewPublisherProcess(name string, cb PublishFunc) (*cony.Publisher, error) {
	pbl, err := client.NewPublisher(name)
	if err != nil {
		return nil, err
	}

	go cb(pbl)

	return pbl, nil
}

// NewPublisher creates new *cony.Publisher to exchange by routing_key.
// NewPublisher does not declare the queue, so if it was not created earlier,
// the publication will go to the void.
func (client *Client) NewPublisher(name string) (*cony.Publisher, error) {
	publisherConfig, ok := client.config.Publishers[name]
	if !ok {
		return nil, ErrEmptyRoutingKeyOrExchangeName
	}

	pbl := cony.NewPublisher(publisherConfig.ExchangeName, publisherConfig.RoutingKey)
	client.cony.Publish(pbl)

	return pbl, nil
}

// NewConsumerBindTo creates new *cony.Consumer and binds to the queue
// by routing_key.
func (client *Client) NewConsumer(name string, opts ...cony.ConsumerOpt) (*cony.Consumer, error) {
	consumerConfig, ok := client.config.Consumers[name]
	if !ok {
		return nil, ErrEmptyRoutingKeyOrQueueName
	}

	var declares []cony.Declaration

	if consumerConfig.RoutingKey != "" {
		que := &cony.Queue{
			Name:       consumerConfig.QueueName,
			AutoDelete: client.config.Server.Queue.AutoDelete,
			Durable:    client.config.Server.Queue.Durable,
			Exclusive:  client.config.Server.Queue.Exclusive,
		}

		declares = append(declares, cony.DeclareQueue(que))

		bnd := cony.Binding{
			Queue: que,
			Exchange: cony.Exchange{
				Name:       client.config.Server.Exchange.Name,
				Kind:       client.config.Server.Exchange.Type,
				AutoDelete: client.config.Server.Exchange.AutoDelete,
				Durable:    client.config.Server.Exchange.Durable,
			},
			Key: consumerConfig.RoutingKey,
		}

		declares = append(declares, cony.DeclareBinding(bnd))
	}

	if len(declares) > 0 {
		client.cony.Declare(declares)
	}

	if client.config.Server.Qos != 0 {
		opts = append(opts, cony.Qos(client.config.Server.Qos))
	}

	cns := cony.NewConsumer(
		&cony.Queue{
			Name:       consumerConfig.QueueName,
			Durable:    client.config.Server.Queue.Durable,
			AutoDelete: client.config.Server.Queue.AutoDelete,
			Exclusive:  client.config.Server.Queue.Exclusive,
		},
		opts...,
	)
	client.cony.Consume(cns)

	return cns, nil
}

func createEnv(cfg *ServerConfig) error {
	client := cony.NewClient(cony.URL(cfg.Server), cony.Backoff(cony.DefaultBackoff))
	defer client.Close()

	var declares []cony.Declaration
	var exc *cony.Exchange
	var que *cony.Queue

	if len(cfg.Exchange.Name) > 0 {
		exc = &cony.Exchange{
			Name:       cfg.Exchange.Name,
			Kind:       cfg.Exchange.Type,
			AutoDelete: cfg.Exchange.AutoDelete,
			Durable:    cfg.Exchange.Durable,
		}

		declares = append(declares, cony.DeclareExchange(*exc))
	}

	if cfg.Queue.Name != "" {
		que = &cony.Queue{
			Name:       cfg.Queue.Name,
			AutoDelete: cfg.Queue.AutoDelete,
			Durable:    cfg.Queue.Durable,
			Exclusive:  cfg.Queue.Exclusive,
			Args:       normalizeArgs(cfg.Queue.Args),
		}

		declares = append(declares, cony.DeclareQueue(que))
	}

	if exc != nil && que != nil {
		bnd := cony.Binding{
			Queue:    que,
			Exchange: *exc,
			Key:      cfg.RoutingKey,
		}

		declares = append(declares, cony.DeclareBinding(bnd))
	}

	if len(declares) > 0 {
		client.Declare(declares)
	}

	client.Loop()
	select {
	case err := <-client.Errors():
		return err
	default:
		return nil
	}
}

func normalizeArgs(args map[string]interface{}) map[string]interface{} {
	newArgs := make(map[string]interface{})

	for key, value := range args {
		switch strings.ToLower(key) {
		case "x-max-priority":
			switch value := value.(type) {
			case int:
				newArgs[key] = uint8(value)
			case int8:
				newArgs[key] = uint8(value)
			case int16:
				newArgs[key] = uint8(value)
			case int32:
				newArgs[key] = uint8(value)
			case int64:
				newArgs[key] = uint8(value)
			case float32:
				newArgs[key] = uint8(value)
			case float64:
				newArgs[key] = uint8(value)
			default:
				newArgs[key] = uint8(0)
			}
		default:
			newArgs[key] = value
		}
	}

	return newArgs
}
