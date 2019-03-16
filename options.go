package broker

import (
	"context"
	"crypto/tls"

	"github.com/rai-project/serializer"
)

// Options ...
type Options struct {
	Endpoints  []string
	Serializer serializer.Serializer
	Secure     bool
	TLSConfig  *tls.Config
	Context    context.Context
}

// Option ...
type Option func(*Options)

// Endpoints ...
func Endpoints(addrs []string) Option {
	return func(o *Options) {
		o.Endpoints = addrs
	}
}

// Serializer ...
func Serializer(s serializer.Serializer) Option {
	return func(o *Options) {
		o.Serializer = s
	}
}

// Secure ...
func Secure(b bool) Option {
	return func(o *Options) {
		o.Secure = b
	}
}

// TLSConfig ...
func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

// PublishOptions ...
type PublishOptions struct {
	Context context.Context
}

// PublishOption ...
type PublishOption func(*PublishOptions)

// SubscribeOptions ...
type SubscribeOptions struct {
	AutoAck                      bool
	Queue                        string
	Context                      context.Context
	BeforeReceiveMessageCallback []func()
	OnReceiveMessageCallback     []func(*Message)
	AfterReceiveMessageCallback  []func()
}

// SubscribeOption ...
type SubscribeOption func(*SubscribeOptions)

// AutoAck ...
func AutoAck(b bool) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.AutoAck = b
	}
}

// Queue ...
func Queue(s string) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Queue = s
	}
}

// BeforeReceiveSubscribeMessageCallback ...
func BeforeReceiveSubscribeMessageCallback(f func()) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.BeforeReceiveMessageCallback = append(o.BeforeReceiveMessageCallback, f)
	}
}

func OnReceiveSubscribeMessageCallback(f func(*Message)) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.OnReceiveMessageCallback = append(o.OnReceiveMessageCallback, f)
	}
}

func AfterReceiveSubscribeMessageCallback(f func()) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.AfterReceiveMessageCallback = append(o.AfterReceiveMessageCallback, f)
	}
}
