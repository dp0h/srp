package pool

import (
	"github.com/dp0h/srp/app/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

const srpYaml1 = `
services:
  srv1:
    host: "example1.com"
    healthcheck: "api/ping"
    weight: 1
`

func TestRandomWeightedPool_NextWithOneElement(t *testing.T) {
	conf := config.NewConf(strings.NewReader(srpYaml1))
	pl := NewRandomWeightedPool(conf, time.Duration(0), time.Duration(0))

	host, err := pl.Next()
	require.Nil(t, err)
	assert.Equal(t, host, "example1.com")
}

const srpYaml2 = `
services:
  srv1:
    host: "example1.com"
    healthcheck: "api/ping"
    weight: 1

  srv2:
    host: "example2.com"
    healthcheck: "api/ping"
    weight: 1

  srv3:
    host: "example3.com"
    healthcheck: "api/ping"
    weight: 0

  srv4:
    host: "example4.com"
    healthcheck: "api/ping"
    weight: 2
`

func TestRandomWeightedPool_NextDistribution(t *testing.T) {
	conf := config.NewConf(strings.NewReader(srpYaml2))
	pl := NewRandomWeightedPool(conf, time.Duration(0), time.Duration(0))

	hosts := make(map[string]int)

	for i := 0; i < 1000; i++ {
		host, err := pl.Next()
		require.Nil(t, err)

		hosts[host]++
	}

	assert.Len(t, hosts, 3)

	assert.Greater(t, hosts["example1.com"], 200)
	assert.Less(t, hosts["example1.com"], 300)

	assert.Greater(t, hosts["example2.com"], 200)
	assert.Less(t, hosts["example2.com"], 300)

	assert.Greater(t, hosts["example4.com"], 400)
	assert.Less(t, hosts["example4.com"], 600)
}
