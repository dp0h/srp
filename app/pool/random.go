package pool

import (
	"errors"
	"github.com/dp0h/srp/app/config"
	"github.com/rs/zerolog/log"
	"math/rand"
	"sync"
	"time"
)

// RandomWeightedPool defines pool of services
type RandomWeightedPool struct {
	refresh  time.Duration
	timeout  time.Duration
	services []*service
	lock     sync.RWMutex
}

type service struct {
	host        string
	healthCheck string
	weight      int
	alive       bool
}

// NewRandomWeightedPool new random weighted pool of services
func NewRandomWeightedPool(config *config.ConfFile, refresh time.Duration, timeout time.Duration) *RandomWeightedPool {
	rand.Seed(int64(time.Now().Nanosecond()))
	res := RandomWeightedPool{
		services: configToServices(config),
		refresh:  refresh,
		timeout:  timeout}
	go res.healthCheck()
	return &res
}

func configToServices(config *config.ConfFile) []*service {
	var res []*service

	for _, v := range config.Services {
		svc := service{
			host:        v.Host,
			healthCheck: v.HealthCheck,
			weight:      v.Weight,
			alive:       true,
		}

		res = append(res, &svc)
	}

	if len(res) == 0 {
		log.Fatal().Msg("no services found")
	}

	return res
}

// Next returns next url
func (p *RandomWeightedPool) Next() (string, error) {
	alive, err := p.getAlive()
	if err != nil {
		return "", err
	}

	if len(alive) == 1 {
		return alive[0].host, nil
	}

	totalWeight := 0
	for _, item := range alive {
		totalWeight += item.weight
	}

	r := rand.Intn(totalWeight) + 1
	for _, item := range alive {
		r -= item.weight
		if r <= 0 {
			return item.host, nil
		}
	}

	return alive[0].host, nil
}

type hostAndWeight struct {
	host   string
	weight int
}

func (p *RandomWeightedPool) getAlive() ([]hostAndWeight, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if len(p.services) == 0 {
		return nil, errors.New("empty services list")
	}

	var res []hostAndWeight

	for _, svc := range p.services {
		if svc.alive && svc.weight > 0 {
			hw := hostAndWeight{
				host:   svc.host,
				weight: svc.weight,
			}
			res = append(res, hw)
		}
	}

	if len(res) == 0 {
		return nil, errors.New("no alive services")
	}

	return res, nil
}
