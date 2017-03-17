package nsq

import (
	"strings"

	"github.com/k0kubun/pp"
	"github.com/rai-project/config"
	"github.com/rai-project/serializer"
	"github.com/rai-project/serializer/bson"
	"github.com/rai-project/serializer/json"
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

var (
	Config = &nsqConfig{
		done: make(chan struct{}),
	}
)

func (nsqConfig) ConfigName() string {
	return "NSQ"
}

func (a *nsqConfig) SetDefaults() {
	vipertags.SetDefaults(a)
}

func (a *nsqConfig) Read() {
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

func (c nsqConfig) Wait() {
	<-c.done
}

func (c nsqConfig) String() string {
	return pp.Sprintln(c)
}

func (c nsqConfig) Debug() {
	log.Debug("NSQ Config = ", c)
}

func init() {
	config.Register(Config)
}
