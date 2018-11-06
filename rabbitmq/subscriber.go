package rabbitmq

import (
	"context"
	"github.com/rai-project/broker"
)

type subscriber struct {
	topic      string
	opts       broker.SubscribeOptions
	cancelFunc context.CancelFunc
}

// Options ...
func (s *subscriber) Options() broker.SubscribeOptions {
	return s.opts
}

// Topic ...
func (s *subscriber) Topic() string {
	return s.topic
}

// Unsubscribe ...
func (s *subscriber) Unsubscribe() error {
	s.cancelFunc()
	return nil
}
