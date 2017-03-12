package sqs

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/rai-project/broker"
)

const (
	sessionKey                = "github.com/rai-project/broker/sqs/session"
	queueNameKey              = "github.com/rai-project/broker/sqs/queueName"
	concurrentHandlerCountKey = "github.com/rai-project/broker/sqs/concurrentHandlerCount"
	subscriptionTimeoutKey    = "github.com/rai-project/broker/sqs/subscriptionTimeout"
	maxInFlightKey            = "github.com/rai-project/broker/sqs/maxInFlight"
)

var (
	DefaultConcurrentHandlerCount = 1
	DefaultSubscriptionTimeout    = int64(5) // second
)

func Session(sess *session.Session) broker.Option {
	return func(o *broker.Options) {
		o.Context = context.WithValue(o.Context, sessionKey, sess)
	}
}

func QueueName(q string) broker.Option {
	return func(o *broker.Options) {
		o.Context = context.WithValue(o.Context, queueNameKey, q)
	}
}

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

func MaxInFlight(n int) broker.Option {
	return func(o *broker.Options) {
		o.Context = context.WithValue(o.Context, maxInFlightKey, n)
	}
}
