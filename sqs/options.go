package sqs

import (
	"context"
	"time"

	"github.com/rai-project/broker"
)

const (
	concurrentHandlerCountKey = "github.com/rai-project/broker/sqs/concurrentHandlerCount"
	subscriptionTimeoutKey    = "github.com/rai-project/broker/sqs/subscriptionTimeout"
)

var (
	DefaultConcurrentHandlerCount = 1
	DefaultSubscriptionTimeout    = int64(5) // second
)

func ConcurrentHandlerCount(n int) broker.SubscribeOption {
	return func(o *broker.SubscribeOptions) {
		o.Context = context.WithValue(o.Context, concurrentHandlerCountKey, n)
	}
}

func SubscriptionTimeout(d time.Duration) broker.SubscribeOption {
	return func(o *broker.SubscribeOptions) {
		o.Context = context.WithValue(o.Context, subscriptionTimeoutKey, int64(d.Seconds()))
	}
}
