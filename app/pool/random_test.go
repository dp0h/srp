package pool

import (
	"github.com/dp0h/srp/app/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

const srpYaml = `
services:
  srv1:
    host: "example1.com"
    healthcheck: "api/ping"
    weight: 1

  srv2:
    host: "example2.com"
    healthcheck: "api/ping"
    weight: 2
`

func TestRandomWeightedPool_Next(t *testing.T) {
	conf := config.NewConf(strings.NewReader(srpYaml))
	pl := NewRandomWeightedPool(conf, time.Duration(0), time.Duration(0))

	host, err := pl.Next()
	require.Nil(t, err)
	assert.Equal(t, host, "example1.com")
}
