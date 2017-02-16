package sqs

import (
	"context"

	"github.com/rai-project/broker"
)

const (
	concurrentHandlerCountKey = "github.com/rai-project/broker/sqs/concurrentHandlerCount"
)

var (
	DefaultConcurrentHandlerCount = 1
)

func ConcurrentHandlerCount(n int) broker.SubscribeOption {
	return func(o *broker.SubscribeOptions) {
		o.Context = context.WithValue(o.Context, concurrentHandlerCountKey, n)
	}
}
