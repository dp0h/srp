package config

import (
	"io"
	"io/ioutil"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

// ConfigFile contains a list service configurations
type ConfFile struct {
	Services map[string]ServiceConfig `yaml:"services"`
}

// ServiceConfig contains single host configuration
type ServiceConfig struct {
	Host        string `yaml:"host"`
	HealthCheck string `yaml:"healthcheck"`
	Weight      int    `yaml:"weight"`
}

// NewConf makes new config for yml reader
func NewConf(reader io.Reader) *ConfFile {
	res := &ConfFile{}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read config")
	}
	if err = yaml.Unmarshal(data, &res); err != nil {
		log.Fatal().Err(err).Msg("failed to parse config")
	}
	if len(res.Services) == 0 {
		log.Fatal().Msg("no valid configuration found")
	}

	return res
}