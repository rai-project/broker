package rabbitmq

import (
	"context"
	"time"

	"github.com/rai-project/broker"
)

const (
	queueNameKey              = "github.com/rai-project/broker/rabbitmq/queueName"
	concurrentHandlerCountKey = "github.com/rai-project/broker/rabbitmq/concurrentHandlerCount"
	subscriptionTimeoutKey    = "github.com/rai-project/broker/rabbitmq/subscriptionTimeout"
	maxInFlightKey            = "github.com/rai-project/broker/rabbitmq/maxInFlight"
)

// DefaultConcurrentHandlerCount ...
var (
	DefaultConcurrentHandlerCount = 1
	DefaultSubscriptionTimeout    = int64(5) // second
)

// QueueName ...
func QueueName(q string) broker.Option {
	return func(o *broker.Options) {
		o.Context = context.WithValue(o.Context, queueNameKey, q)
	}
}

// ConcurrentHandlerCount ...
func ConcurrentHandlerCount(n int) broker.SubscribeOption {
	return func(o *broker.SubscribeOptions) {
		o.Context = context.WithValue(o.Context, concurrentHandlerCountKey, n)
	}
}

// SubscriptionTimeout ...
func SubscriptionTimeout(d time.Duration) broker.SubscribeOption {
	return func(o *broker.SubscribeOptions) {
		o.Context = context.WithValue(o.Context, subscriptionTimeoutKey, int64(d.Seconds()))
	}
}

// MaxInFlight ...
func MaxInFlight(n int) broker.Option {
	return func(o *broker.Options) {
		o.Context = context.WithValue(o.Context, maxInFlightKey, n)
	}
}
