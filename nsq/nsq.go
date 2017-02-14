package nsq

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"math/rand"
	"sync"

	"github.com/k0kubun/pp"
	nsq "github.com/nsqio/go-nsq"
	"github.com/pkg/errors"
	"github.com/rai-project/broker"
	"github.com/spf13/viper"
)

type nsqBroker struct {
	sync.Mutex
	isRunning   bool
	publishers  []*publisher
	subscribers []*subscriber
	opts        broker.Options
	config      *nsq.Config
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
		Endpoints:  Config.NsqdEndpoints,
		Serializer: Config.Serializer,
		Secure:     Config.CACertificate != "",
		TLSConfig:  tlsConf,
		Context:    context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	config := nsq.NewConfig()

	if viper.IsSet("app.name") {
		config.UserAgent = fmt.Sprintf("%s/go-nsq", viper.GetString("app.name"))
	}

	if maxInFlight, ok := options.Context.Value(maxInFlightKey).(int); ok {
		config.MaxInFlight = maxInFlight
	}

	return &nsqBroker{
		isRunning:   false,
		opts:        options,
		config:      config,
		publishers:  []*publisher{},
		subscribers: []*subscriber{},
	}
}

func (b *nsqBroker) Options() broker.Options {
	b.Lock()
	defer b.Unlock()
	return b.opts
}

func (b *nsqBroker) Connect() error {
	b.Lock()
	defer b.Unlock()

	if b.isRunning {
		return nil
	}

	publishers := []*publisher{}
	for _, addr := range b.Options().Endpoints {
		p, err := nsq.NewProducer(addr, b.config)
		if err != nil {
			return errors.Wrapf(err,
				"Failed to create new producer at address = %v", addr)
		}
		publishers = append(publishers, &publisher{p})
	}

	for _, subscriber := range b.subscribers {
		err := subscriber.c.ConnectToNSQDs(b.Options().Endpoints)
		if err != nil {
			return errors.Wrapf(err,
				"Failed to connect to nsqd at address = %v", pp.Sprint(b.Options().Endpoints))
		}
	}

	return nil
}

func (b *nsqBroker) Disconnect() error {
	b.Lock()
	defer b.Unlock()

	for _, publisher := range b.publishers {
		publisher.Stop()
	}

	for _, subscriber := range b.subscribers {
		subscriber.c.Stop()
		for _, addr := range b.Options().Endpoints {
			subscriber.c.DisconnectFromNSQD(addr)
		}
	}

	b.publishers = []*publisher{}
	b.subscribers = []*subscriber{}
	b.isRunning = false

	return nil

}

func (b *nsqBroker) Publish(queue string, msg *broker.Message, opts ...broker.PublishOption) error {
	lp := len(b.publishers)
	if lp == 0 {
		return errors.New("No publishers found")
	}
	p := b.publishers[rand.Intn(lp)]
	bts, err := b.opts.Serializer.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "Failed to serialize message while trying to publish to queue")
	}
	return p.Publish(queue, bts)
}

func (b *nsqBroker) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {

	options := broker.SubscribeOptions{
		AutoAck: Config.AutoAck,
		Queue:   "",
		Context: context.WithValue(
			context.Background(),
			concurrentHandlerCountKey,
			DefaultConcurrentHandlerCount,
		),
	}

	for _, o := range opts {
		o(&options)
	}

	channel := "#ch"
	if options.Queue != "" {
		channel = options.Queue
	}

	consumer, err := nsq.NewConsumer(topic, channel, b.config)
	if err != nil {
		return nil, errors.Wrapf(err,
			"Failed to create new consumer on topic=%v and channel=%v", topic, channel)
	}

	nsqHandler := nsq.HandlerFunc(func(nm *nsq.Message) error {
		if !options.AutoAck {
			nm.DisableAutoResponse()
		}
		msg := new(broker.Message)
		err := b.opts.Serializer.Unmarshal(nm.Body, msg)
		if err != nil {
			return errors.Wrap(err,
				"Failed to unmarshal message while handeling queue message")
		}
		return handler(&publication{
			topic:   topic,
			channel: channel,
			m:       msg,
			nm:      nm,
		})
	})

	concurrentHandlerCount, ok := options.Context.Value(concurrentHandlerCountKey).(int)
	if ok {
		consumer.AddConcurrentHandlers(nsqHandler, concurrentHandlerCount)
	}

	return &subscriber{
		c:                      consumer,
		topic:                  topic,
		channel:                channel,
		handler:                nsqHandler,
		concurrentHandlerCount: concurrentHandlerCount,
		opts: options,
	}, nil
}

func (b *nsqBroker) Name() string {
	return "nsq"
}
