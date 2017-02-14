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
