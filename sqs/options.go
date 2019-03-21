package sqs

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/rai-project/broker"
)

type sessionKey struct{}
type queueNameKey struct{}
type concurrentHandlerCountKey struct{}
type subscriptionTimeoutKey struct{}
type maxInFlightKey struct{}

// DefaultConcurrentHandlerCount ...
var (
	DefaultConcurrentHandlerCount = 1
	DefaultSubscriptionTimeout    = int64(5) // second
)

// Session ...
func Session(sess *session.Session) broker.Option {
	return func(o *broker.Options) {
		o.Context = context.WithValue(o.Context, sessionKey{}, sess)
	}
}

// QueueName ...
func QueueName(q string) broker.Option {
	return func(o *broker.Options) {
		o.Context = context.WithValue(o.Context, queueNameKey{}, q)
	}
}

// ConcurrentHandlerCount ...
func ConcurrentHandlerCount(n int) broker.SubscribeOption {
	return func(o *broker.SubscribeOptions) {
		o.Context = context.WithValue(o.Context, concurrentHandlerCountKey{}, n)
	}
}

// SubscriptionTimeout ...
func SubscriptionTimeout(d time.Duration) broker.SubscribeOption {
	return func(o *broker.SubscribeOptions) {
		o.Context = context.WithValue(o.Context, subscriptionTimeoutKey{}, int64(d.Seconds()))
	}
}

// MaxInFlight ...
func MaxInFlight(n int) broker.Option {
	return func(o *broker.Options) {
		o.Context = context.WithValue(o.Context, maxInFlightKey{}, n)
	}
}
