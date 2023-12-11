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
    host: "http://example1.com"
    healthcheck: "api/ping"
    weight: 1
`

func TestRandomWeightedPool_NextWithOneElement(t *testing.T) {
	conf := config.NewConf(strings.NewReader(srpYaml1))
	pl := NewRandomWeightedPool(conf, time.Duration(0), time.Duration(0))

	host, err := pl.Next()
	require.Nil(t, err)
	assert.Equal(t, host, "http://example1.com")
}

const srpYaml2 = `
services:
  srv1:
    host: "http://example1.com"
    healthcheck: "api/ping"
    weight: 1

  srv2:
    host: "http://example2.com"
    healthcheck: "api/ping"
    weight: 1

  srv3:
    host: "http://example3.com"
    healthcheck: "api/ping"
    weight: 0

  srv4:
    host: "http://example4.com"
    healthcheck: "api/ping"
    weight: 2

  srv5:
    host: "http://example5.com"
    healthcheck: "api/ping"
    weight: 2
`

func TestRandomWeightedPool_NextDistribution(t *testing.T) {
	conf := config.NewConf(strings.NewReader(srpYaml2))
	pl := NewRandomWeightedPool(conf, time.Duration(0), time.Duration(0))

	// set no alive
	for _, item := range pl.services {
		if item.host == "http://example4.com" {
			item.alive = false
		}
	}

	hosts := make(map[string]int)

	for i := 0; i < 1000; i++ {
		host, err := pl.Next()
		require.Nil(t, err)

		hosts[host]++
	}

	assert.Len(t, hosts, 3)

	assert.Greater(t, hosts["http://example1.com"], 200)
	assert.Less(t, hosts["http://example1.com"], 300)

	assert.Greater(t, hosts["http://example2.com"], 200)
	assert.Less(t, hosts["http://example2.com"], 300)

	assert.Greater(t, hosts["http://example5.com"], 400)
	assert.Less(t, hosts["http://example5.com"], 600)
}
