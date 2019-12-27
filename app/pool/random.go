package pool

import (
	"errors"
	"github.com/dp0h/srp/app/config"
	"github.com/rs/zerolog/log"
	"math/rand"
	"sync"
	"time"
)

// RandomWeighted
type RandomWeightedPool struct {
	refresh  time.Duration
	timeout  time.Duration
	services []*service
	lock     sync.RWMutex
}

type service struct {
	Host        string
	HealthCheck string
	Weight      int
	alive       bool
}

// NewRandomWeighted new random weighted pool of services
func NewRandomWeightedPool(config *config.ConfFile, refresh time.Duration, timeout time.Duration) *RandomWeightedPool {
	rand.Seed(int64(time.Now().Nanosecond()))
	res := RandomWeightedPool{
		services: configToServices(config),
		refresh:  refresh,
		timeout:  timeout}
	return &res
}

func configToServices(config *config.ConfFile) []*service {
	var res []*service

	for _, v := range config.Services {
		svc := service{
			Host:        v.Host,
			HealthCheck: v.HealthCheck,
			Weight:      v.Weight,
			alive:       true,
		}

		res = append(res, &svc)
	}

	if len(res) == 0 {
		log.Fatal().Msg("no services found")
	}

	return res
}

func (p *RandomWeightedPool) Next() (string, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if len(p.services) == 0 {
		return "", errors.New("empty services list")
	}

	return p.services[0].Host, nil
}
