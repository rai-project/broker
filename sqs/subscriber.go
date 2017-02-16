package sqs

import (
	"context"

	"github.com/rai-project/broker"
)

type subscriber struct {
	topic      string
	cancelFunc context.CancelFunc
	opts       broker.SubscribeOptions
}

func (s *subscriber) Options() broker.SubscribeOptions {
	return s.opts
}

func (s *subscriber) Topic() string {
	return s.topic
}

func (s *subscriber) Unsubscribe() error {
	s.cancelFunc()
	return nil
}
