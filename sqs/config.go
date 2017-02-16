package sqs

import (
	"strings"

	"github.com/k0kubun/pp"
	"github.com/rai-project/config"
	"github.com/rai-project/serializer"
	"github.com/rai-project/serializer/bson"
	"github.com/rai-project/serializer/json"
	"github.com/rai-project/vipertags"
)

type sqsConfig struct {
	Provider       string                `json:"provider" config:"broker.provider" default:"sqs"`
	Serializer     serializer.Serializer `json:"-" config:"-"`
	SerializerName string                `json:"serializer_name" config:"broker.serializer" default:"json"`
	Endpoints      []string              `json:"endpoints" config:"broker.endpoints"`
	AutoAck        bool                  `json:"autoack" config:"broker.autoack" default:"true"`
}

var (
	Config = &sqsConfig{}
)

func (sqsConfig) ConfigName() string {
	return "SQS"
}

func (sqsConfig) SetDefaults() {
}

func (a *sqsConfig) Read() {
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

func (c sqsConfig) String() string {
	return pp.Sprintln(c)
}

func (c sqsConfig) Debug() {
	log.Debug("SQS Config = ", c)
}

func init() {
	config.Register(Config)
}
