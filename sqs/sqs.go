package sqs

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
	raiaws "github.com/rai-project/aws"
	"github.com/rai-project/broker"
)

type sqsBroker struct {
	sync.Mutex
	isRunning bool
	session   *session.Session
	svc       *sqs.SQS
	opts      broker.Options
}

func New(opts ...broker.Option) broker.Broker {

	options := broker.Options{
		Endpoints:  Config.Endpoints,
		Serializer: Config.Serializer,
		Context:    context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	return &sqsBroker{
		isRunning: false,
		session:   nil,
		svc:       nil,
		opts:      options,
	}
}

func (b *sqsBroker) Options() broker.Options {
	b.Lock()
	defer b.Unlock()
	return b.opts
}

func (b *sqsBroker) Connect() error {
	b.Lock()
	defer b.Unlock()

	if b.isRunning {
		return nil
	}

	// Initialize a session that the SDK will use to load configuration,
	// credentials, and region from the shared config file. (~/.aws/config).
	sess, err := raiaws.NewSession()
	if err != nil {
		return err
	}
	b.session = sess

	// Create a SQS service client.
	if len(b.opts.Endpoints) == 0 {
		return errors.New("Invalid sqs endpoint")
	}
	endpoint := b.opts.Endpoints[0]
	svc := sqs.New(sess, aws.NewConfig().WithEndpoint(endpoint))
	b.svc = svc

	b.isRunning = true

	return nil
}

func (b *sqsBroker) Disconnect() error {
	b.Lock()
	defer b.Unlock()

	b.isRunning = false

	return nil

}

func (b *sqsBroker) Publish(queue string, msg *broker.Message, opts ...broker.PublishOption) error {
	bts, err := b.opts.Serializer.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "Failed to serialize message while trying to publish to queue")
	}
	resultURL, err := b.svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queue),
	})
	if err != nil {
		return errors.Wrapf(err, "Unable to queue %v.", queue)
	}

	_, err = b.svc.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    resultURL.QueueUrl,
		MessageBody: aws.String(string(bts)),
	})
	if err != nil {
		return errors.Wrapf(err, "Failed to send message to sqs.")
	}
	return nil
}

func (b *sqsBroker) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {

	options := broker.SubscribeOptions{
		AutoAck: Config.AutoAck,
		Queue:   "",
		Context: context.WithValue(
			context.Background(),
			concurrentHandlerCountKey,
			DefaultConcurrentHandlerCount,
		),
	}

	for _, o := range opts {
		o(&options)
	}

	timeout := int64(1) // Second

	concurrentHandlerCount0, ok := options.Context.Value(concurrentHandlerCountKey).(int)
	if !ok {
		concurrentHandlerCount0 = 1
	}
	concurrentHandlerCount := int64(concurrentHandlerCount0)

	resultURL, err := b.svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(topic),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to queue %v.", topic)
	}

	cancelCtx, cancelFunc := context.WithCancel(options.Context)

	go func() {
		for {
			select {
			case <-cancelCtx.Done():
				return
			default:
				// Receive a message from the SQS queue with long polling enabled.
				result, err := b.svc.ReceiveMessage(&sqs.ReceiveMessageInput{
					QueueUrl: resultURL.QueueUrl,
					AttributeNames: aws.StringSlice([]string{
						"SentTimestamp",
					}),
					MaxNumberOfMessages: aws.Int64(concurrentHandlerCount),
					MessageAttributeNames: aws.StringSlice([]string{
						"All",
					}),
					WaitTimeSeconds: aws.Int64(timeout),
				})
				if err != nil {
					log.Errorf("Unable to receive message from queue %q, %v.", topic, err)
					continue
				}

				if !options.AutoAck {
					for _, msg := range result.Messages {
						b.svc.DeleteMessage(&sqs.DeleteMessageInput{
							QueueUrl:      resultURL.QueueUrl,
							ReceiptHandle: msg.ReceiptHandle,
						})
					}
				}

				for _, m := range result.Messages {
					msg := new(broker.Message)
					err := b.opts.Serializer.Unmarshal([]byte(*m.Body), msg)
					if err != nil {
						log.WithError(err).
							Errorf("Failed to unmarshal message while handeling queue message")
					}

					handler(&publication{
						svc:           b.svc,
						topic:         topic,
						m:             msg,
						nm:            m,
						queueUrl:      *resultURL.QueueUrl,
						receiptHandle: *m.ReceiptHandle,
					})
				}
			}
		}
	}()

	return &subscriber{
		topic:      topic,
		opts:       options,
		cancelFunc: cancelFunc,
	}, nil
}

func (b *sqsBroker) Name() string {
	return "sqs"
}
