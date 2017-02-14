package rabbitmq

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"sync"

	"github.com/rai-project/broker"
	"github.com/streadway/amqp"
)

type rabbitmqBroker struct {
	sync.Mutex
	isRunning   bool
	publishers  []*publisher
	subscribers []*subscriber
	opts        broker.Options
	config      amqp.Config
	conn        *amqp.Connection
}

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
		Endpoints:  Config.Endpoints,
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

	return &nsqBroker{
		isRunning:   false,
		opts:        options,
		config:      config,
		publishers:  []*publisher{},
		subscribers: []*subscriber{},
	}
}

func (b *rabbitmqBroker) Options() broker.Options {
	b.Lock()
	defer b.Unlock()
	return b.opts
}

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

	return nil
}

func (b *rabbitmqBroker) Disconnect() error {
	b.Lock()
	defer b.Unlock()

	if err := b.conn.Close(); err != nil {
		return err
	}

	b.isRunning = false

	return
}

func (b *rabbitmqBroker) Name() string {
	return "rabbitmq"
}
