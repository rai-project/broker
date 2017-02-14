package broker

import (
	"context"
	"crypto/tls"

	"github.com/rai-project/serializer"
)

type Options struct {
	Endpoints  []string
	Serializer serializer.Serializer
	Secure     bool
	TLSConfig  *tls.Config
	Context    context.Context
}

type Option func(*Options)

func Endpoints(addrs []string) Option {
	return func(o *Options) {
		o.Endpoints = addrs
	}
}

func Serializer(s serializer.Serializer) Option {
	return func(o *Options) {
		o.Serializer = s
	}
}

func Secure(b bool) Option {
	return func(o *Options) {
		o.Secure = b
	}
}

func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

type PublishOptions struct {
	Context context.Context
}

type PublishOption func(*PublishOptions)

type SubscribeOptions struct {
	AutoAck bool
	Queue   string
	Context context.Context
}

type SubscribeOption func(*SubscribeOptions)

func AutoAck(b bool) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.AutoAck = b
	}
}

func Queue(s string) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Queue = s
	}
}
