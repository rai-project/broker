package rabbitmq


import (
	"github.com/rai-project/broker"
)

type publication struct {
	topic         string
	m             *broker.Message
	queue         string
	opts          broker.PublishOptions
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
	return nil
}
