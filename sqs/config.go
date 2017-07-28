package sqs

import (
	"github.com/k0kubun/pp"
	"github.com/rai-project/config"
	"github.com/rai-project/serializer"
	_ "github.com/rai-project/serializer/bson"
	_ "github.com/rai-project/serializer/json"
	"github.com/rai-project/vipertags"
)

type sqsConfig struct {
	Provider       string                `json:"provider" config:"broker.provider"`
	Serializer     serializer.Serializer `json:"-" config:"-"`
	SerializerName string                `json:"serializer_name" config:"broker.serializer" default:"json"`
	AutoAck        bool                  `json:"autoack" config:"broker.autoack" default:"true"`
	Region         string                `json:"region" config:"broker.region" default:"us-east-1"`
	done           chan struct{}         `json:"-" config:"-"`
}

var (
	Config = &sqsConfig{
		done: make(chan struct{}),
	}
)

func (sqsConfig) ConfigName() string {
	return "SQS"
}

func (a *sqsConfig) SetDefaults() {
	vipertags.SetDefaults(a)
}

func (a *sqsConfig) Read() {
	defer close(a.done)
	vipertags.Fill(a)
	a.Serializer, _ = serializer.FromName(a.SerializerName)
}

func (c sqsConfig) Wait() {
	<-c.done
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
