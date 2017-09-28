package nsq

import (
	"github.com/rai-project/broker"
	"github.com/rai-project/config"
	"github.com/rai-project/logger"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Entry
)

func init() {
	config.AfterInit(func() {
		log = logger.New().WithField("pkg", "nsq")
		broker.Standard = New()
	})
}
