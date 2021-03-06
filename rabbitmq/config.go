package rabbitmq

import (
	"github.com/k0kubun/pp"
	"github.com/rai-project/config"
	"github.com/rai-project/serializer"
	_ "github.com/rai-project/serializer/bson"
	_ "github.com/rai-project/serializer/json"
	_ "github.com/rai-project/serializer/jsonpb"
	"github.com/rai-project/vipertags"
)

type rabbitmqConfig struct {
	Provider            string                `json:"provider" config:"broker.provider"`
	Serializer          serializer.Serializer `json:"-" config:"-"`
	SerializerName      string                `json:"serializer_name" config:"broker.serializer" default:"json"`
	RabbitMQEndpoints   []string              `json:"rabbitmq_endpoints" config:"broker.rabbitmq_endpoints"`
	CACertificate       string                `json:"ca_certificate" config:"broker.ca_certificate"`
	AutoAck             bool                  `json:"autoack" config:"broker.autoack" default:"true"`
	done                chan struct{}         `json:"-" config:"-"`
}

// Config ...
var (
	Config = &rabbitmqConfig{
		done: make(chan struct{}),
	}
)

// ConfigName ...
func (rabbitmqConfig) ConfigName() string {
	return "RabbitMQ"
}

// SetDefaults ...
func (a *rabbitmqConfig) SetDefaults() {
	vipertags.SetDefaults(a)
}

// Read ...
func (a *rabbitmqConfig) Read() {
	defer close(a.done)
	vipertags.Fill(a)
	a.Serializer, _ = serializer.FromName(a.SerializerName)
}

// Wait ...
func (c rabbitmqConfig) Wait() {
	<-c.done
}

// String ...
func (c rabbitmqConfig) String() string {
	return pp.Sprintln(c)
}

// Debug ...
func (c rabbitmqConfig) Debug() {
	log.Debug("RabbitMQ Config = ", c)
}

func init() {
	config.Register(Config)
}
