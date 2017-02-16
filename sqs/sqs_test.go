package sqs

import (
	"testing"

	"github.com/rai-project/config"

	"github.com/stretchr/testify/assert"
)

func TestSQS(t *testing.T) {
	config.Init()
	assert.NoError(t, test())
}
