package nsq

import nsq "github.com/nsqio/go-nsq"

type publisher struct {
	*nsq.Producer
}
