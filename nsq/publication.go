package nsq

import (
	nsq "github.com/nsqio/go-nsq"
	"github.com/rai-project/broker"
)

type publication struct {
	topic   string
	channel string
	m       *broker.Message
	nm      *nsq.Message
	opts    broker.PublishOptions
}

// Topic ...
func (p *publication) Topic() string {
	return p.topic
}

// Message ...
func (p *publication) Message() *broker.Message {
	return p.m
}

// Ack ...
func (p *publication) Ack() error {
	p.nm.Finish()
	return nil
}
