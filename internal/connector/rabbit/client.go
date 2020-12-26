// TODO: Try to simplify.
package rabbit

import (
	"fmt"
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

// Config returns a pointer to the Config with which the connection was made.
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

// NewPublisher creates new *cony.Publisher.
func (client *Client) NewPublisher(cb PublishFunc) *cony.Publisher {
	pbl := cony.NewPublisher(client.config.Server.Exchange.Name, client.config.Server.RoutingKey)
	client.cony.Publish(pbl)

	go cb(pbl)

	return pbl
}

// NewPublisher creates new *cony.Publisher and binds to the queue.
func (client *Client) NewPublisherBind(cb PublishFunc) (*cony.Publisher, error) {
	if client.config.Server.RoutingKey == "" {
		return nil, ErrEmptyRoutingKey
	}

	var declares []cony.Declaration

	if client.config.Server.Queue.Name != "" {
		que := &cony.Queue{
			Name:       client.config.Server.Queue.Name,
			AutoDelete: client.config.Server.Queue.AutoDelete,
			Durable:    client.config.Server.Queue.Durable,
			Exclusive:  client.config.Server.Queue.Exclusive,
			Args:       normalizeArgs(client.config.Server.Queue.Args),
		}

		declares = append(declares, cony.DeclareQueue(que))

		bnd := cony.Binding{
			Queue: que, // queue
			Exchange: cony.Exchange{
				Name:       client.config.Server.Exchange.Name,
				Kind:       client.config.Server.Exchange.Type,
				AutoDelete: client.config.Server.Exchange.AutoDelete,
				Durable:    client.config.Server.Exchange.Durable,
			},
			Key: client.config.Server.RoutingKey,
		}

		declares = append(declares, cony.DeclareBinding(bnd))
	}

	if len(declares) > 0 {
		client.cony.Declare(declares)
	}

	pbl := cony.NewPublisher(client.config.Server.Exchange.Name, client.config.Server.RoutingKey)
	client.cony.Publish(pbl)

	go cb(pbl)

	return pbl, nil
}

// NewConsumer creates new *cony.Consumer.
func (client *Client) NewConsumer(opts ...cony.ConsumerOpt) *cony.Consumer {
	if client.config.Server.Qos != 0 {
		opts = append(opts, cony.Qos(client.config.Server.Qos))
	}

	cns := cony.NewConsumer(
		&cony.Queue{
			Name:       client.config.Server.Queue.Name,
			Durable:    client.config.Server.Queue.Durable,
			AutoDelete: client.config.Server.Queue.AutoDelete,
			Exclusive:  client.config.Server.Queue.Exclusive,
		},
		opts...,
	)
	client.cony.Consume(cns)

	return cns
}

// NewConsumerBindTo creates new *cony.Consumer and binds to the queue
// by routing_key.
func (client *Client) NewConsumerBindTo(name string, opts ...cony.ConsumerOpt) (*cony.Consumer, error) {
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
				fmt.Printf("%v %t", args[key], args[key])

				newArgs[key] = uint8(0)
			}
		default:
			newArgs[key] = value
		}
	}

	return newArgs
}
