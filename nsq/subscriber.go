package nsq

import (
	nsq "github.com/nsqio/go-nsq"
	"github.com/rai-project/broker"
)

type subscriber struct {
	c                      *nsq.Consumer
	topic                  string
	channel                string
	handler                nsq.HandlerFunc
	concurrentHandlerCount int
	opts                   broker.SubscribeOptions
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
	s.c.Stop()
	return nil
}
