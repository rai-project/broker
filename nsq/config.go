package nsq

import (
	"github.com/k0kubun/pp"
	"github.com/rai-project/config"
	"github.com/rai-project/serializer"
	_ "github.com/rai-project/serializer/bson"
	_ "github.com/rai-project/serializer/json"
	"github.com/rai-project/vipertags"
)

type nsqConfig struct {
	Provider            string                `json:"provider" config:"broker.provider"`
	Serializer          serializer.Serializer `json:"-" config:"-"`
	SerializerName      string                `json:"serializer_name" config:"broker.serializer" default:"json"`
	NsqdEndpoints       []string              `json:"nsqd_endpoints" config:"broker.nsqd_endpoints"`
	NsqLookupdEndpoints []string              `json:"nsqlookupd_endpoints" config:"broker.nsqlookupd_endpoints"`
	CACertificate       string                `json:"ca_certificate" config:"broker.ca_certificate"`
	ConcurrentHandlers  int                   `json:"concurrent_handlers" config:"broker.concurrent_handlers" default:"1"`
	AutoAck             bool                  `json:"autoack" config:"broker.autoack" default:"true"`
	Ephemeral           bool                  `json:"ephemeral" config:"broker.ephemeral" default:"false"`
	done                chan struct{}         `json:"-" config:"-"`
}

// Config ...
var (
	Config = &nsqConfig{
		done: make(chan struct{}),
	}
)

// ConfigName ...
func (nsqConfig) ConfigName() string {
	return "NSQ"
}

// SetDefaults ...
func (a *nsqConfig) SetDefaults() {
	vipertags.SetDefaults(a)
}

// Read ...
func (a *nsqConfig) Read() {
	defer close(a.done)
	vipertags.Fill(a)
	a.Serializer, _ = serializer.FromName(a.SerializerName)
}

// Wait ...
func (c nsqConfig) Wait() {
	<-c.done
}

// String ...
func (c nsqConfig) String() string {
	return pp.Sprintln(c)
}

// Debug ...
func (c nsqConfig) Debug() {
	log.Debug("NSQ Config = ", c)
}

func init() {
	config.Register(Config)
}
