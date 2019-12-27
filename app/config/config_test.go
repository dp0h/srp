package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestConfig(t *testing.T) {
	conf := NewConf(strings.NewReader(srpYaml))

	assert.NotNil(t, conf)
	assert.Len(t, conf.Services, 2)

	assert.Equal(t, conf.Services["srv1"].Host, "example1.com")
	assert.Equal(t, conf.Services["srv1"].HealthCheck, "api/ping")
	assert.Equal(t, conf.Services["srv1"].Weight, 1)

	assert.Equal(t, conf.Services["srv2"].Host, "example2.com")
	assert.Equal(t, conf.Services["srv2"].HealthCheck, "api/ping")
	assert.Equal(t, conf.Services["srv2"].Weight, 2)
}
