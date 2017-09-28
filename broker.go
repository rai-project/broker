package broker

// Broker ...
type Broker interface {
	Options() Options
	Connect() error
	Disconnect() error
	Publish(queue string, msg *Message, opts ...PublishOption) error
	Subscribe(queue string, handler Handler, opts ...SubscribeOption) (Subscriber, error)
	Name() string
}

// Handler ...
type Handler func(Publication) error

// Message ...
type Message struct {
	ID     string
	Header map[string]string
	Body   []byte
}

// Publication ...
type Publication interface {
	Topic() string
	Message() *Message
	Ack() error
}

// Subscriber ...
type Subscriber interface {
	Options() SubscribeOptions
	Topic() string
	Unsubscribe() error
}
