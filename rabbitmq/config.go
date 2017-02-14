package rabbitmq

import (
	"strings"

	"github.com/k0kubun/pp"
	"github.com/rai-project/config"
	"github.com/rai-project/serializer"
	"github.com/rai-project/serializer/bson"
	"github.com/rai-project/serializer/json"
	"github.com/rai-project/vipertags"
)

type rabbitmqConfig struct {
	Provider            string                `json:"provider" config:"broker.provider" default:"rabbitmq"`
	Serializer          serializer.Serializer `json:"-" config:"-"`
	SerializerName      string                `json:"serializer_name" config:"broker.serializer" default:"json"`
	Endpoints           []string              `json:"endpoints" config:"broker.endpoints"`
	NsqLookupdEndpoints []string              `json:"nsqlookupd_endpoints" config:"broker.nsqlookupd_endpoints"`
	CACertificate       string                `json:"ca_certificate" config:"broker.ca_certificate"`
	AutoAck             bool                  `json:"autoack" config:"broker.autoack" default:"true"`
}

var (
	Config = &rabbitmqConfig{}
)

func (rabbitmqConfig) ConfigName() string {
	return "RabbitMQ"
}

func (rabbitmqConfig) SetDefaults() {
}

func (a *rabbitmqConfig) Read() {
	vipertags.Fill(a)
	switch strings.ToLower(a.SerializerName) {
	case "json":
		a.Serializer = json.New()
	case "bson":
		a.Serializer = bson.New()
	default:
		log.WithField("serializer", a.SerializerName).
			Warn("Cannot find serializer")
	}
}

func (c rabbitmqConfig) String() string {
	return pp.Sprintln(c)
}

func (c rabbitmqConfig) Debug() {
	log.Debug("RabbitMQ Config = ", c)
}

func init() {
	config.Register(Config)
}
