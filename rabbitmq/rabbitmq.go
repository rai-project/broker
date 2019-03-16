package rabbitmq

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"sync"

	"github.com/pkg/errors"
	"github.com/rai-project/broker"
	ctx "github.com/rai-project/context"
	"github.com/streadway/amqp"
)

type rabbitmqBroker struct {
	sync.Mutex
	isRunning bool
	opts      broker.Options
	config    amqp.Config
	conn      *amqp.Connection
}

// New ...
func New(opts ...broker.Option) broker.Broker {

	var tlsConf *tls.Config = nil
	if Config.CACertificate != "" {
		roots := x509.NewCertPool()
		roots.AppendCertsFromPEM([]byte(Config.CACertificate))
		tlsConf = &tls.Config{
			RootCAs: roots,
		}
	}

	options := broker.Options{
		Endpoints:  Config.RabbitMQEndpoints,
		Serializer: Config.Serializer,
		Secure:     Config.CACertificate != "",
		TLSConfig:  tlsConf,
		Context:    context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	config := amqp.Config{
		TLSClientConfig: tlsConf,
	}

	return &rabbitmqBroker{
		isRunning: false,
		opts:      options,
		config:    config,
	}
}

// Options ...
func (b *rabbitmqBroker) Options() broker.Options {
	b.Lock()
	defer b.Unlock()
	return b.opts
}

// Connect ...
func (b *rabbitmqBroker) Connect() error {
	b.Lock()
	defer b.Unlock()

	if b.isRunning {
		return nil
	}

	conn, err := amqp.DialConfig(b.opts.Endpoints[0], b.config)
	if err != nil {
		return err
	}

	b.conn = conn

	b.isRunning = true

	return nil
}

// Disconnect ...
func (b *rabbitmqBroker) Disconnect() error {
	b.Lock()
	defer b.Unlock()

	if err := b.conn.Close(); err != nil {
		return err
	}

	b.isRunning = false

	return nil
}

// Publish ...
func (b *rabbitmqBroker) Publish(queue string, msg *broker.Message, opts ...broker.PublishOption) error {
	bts, err := b.opts.Serializer.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "Failed to serialize message while trying to publish to queue")
	}

	ch, err := b.conn.Channel()
	if err != nil {
		return errors.Wrap(err, "Failed to open a channel")
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return errors.Wrap(err, "Failed to declare a queue")
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(bts),
		})
	if err != nil {
		return errors.Wrap(err, "Failed to send message to RabbitMQ")
	}
	return nil
}

// Subscribe ...
func (b *rabbitmqBroker) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	options := broker.SubscribeOptions{
		AutoAck: Config.AutoAck,
		Queue:   "",
		Context: ctx.Background().
			WithValue(
				concurrentHandlerCountKey,
				DefaultConcurrentHandlerCount,
			).
			WithValue(
				subscriptionTimeoutKey,
				DefaultSubscriptionTimeout,
			),
	}

	for _, o := range opts {
		o(&options)
	}

	cancelCtx, cancelFunc := context.WithCancel(options.Context)

	go func() {
		for {
			select {
			case <-cancelCtx.Done():
				return
			default:
				for _, f := range options.BeforeReceiveMessageCallback {
					f()
				}
				ch, err := b.conn.Channel()
				if err != nil {
					log.WithError(err).Errorf("Failed to open a channel")
				}
				defer ch.Close()

				queueName, ok := b.opts.Context.Value(queueNameKey).(string)
				if !ok {
					log.WithError(err).Errorf("cannot find queue name. make sure to set it when initializing RabbitMQ")
				}

				q, err := ch.QueueDeclare(
					queueName, // name
					true,      // durable
					false,     // delete when unused
					false,     // exclusive
					false,     // no-wait
					nil,       // arguments
				)
				if err != nil {
					log.WithError(err).Errorf("Failed to declare a queue")
				}

				msgs, err := ch.Consume(
					q.Name,          // queue
					"",              // consumer
					options.AutoAck, // auto-ack
					false,           // exclusive
					false,           // no-local
					false,           // no-wait
					nil,             // args
				)
				if err != nil {
					log.WithError(err).Errorf("Failed to register a consumer")
				}
				for m := range msgs {
					msg := new(broker.Message)
					err := b.opts.Serializer.Unmarshal(m.Body, msg)
					if err != nil {
						log.WithError(err).Errorf("Failed to unmarshal message while handeling queue message")
					}

					for _, f := range options.OnReceiveMessageCallback {
						f(msg)
					}

					handler(&publication{
						topic: topic,
						m:     msg,
						queue: queueName,
					})
				}

				for _, f := range options.AfterReceiveMessageCallback {
					f()
				}
			}
		}
	}()
	return &subscriber{
		topic:      topic,
		opts:       options,
		cancelFunc: cancelFunc,
	}, nil
}

// Name ...
func (b *rabbitmqBroker) Name() string {
	return "rabbitmq"
}
