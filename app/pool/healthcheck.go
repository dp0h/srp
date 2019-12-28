package pool

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

// healthCheck periodically checks health status
func (p *RandomWeightedPool) healthCheck() {
	log.Info().Dur("refresh", p.refresh).Dur("timeout", p.timeout).Msg("starting health checks")

	for {
		errorCnt := 0
		for _, item := range p.services {
			if !p.checkService(item) {
				errorCnt++
			}
		}
		log.Debug().Int("total", len(p.services)).Int("failed", errorCnt).Msg("healthcheck update")

		time.Sleep(p.refresh)
	}
}

func (p *RandomWeightedPool) checkService(svc *service) bool {
	if svc.healthCheck == "" {
		return true
	}

	err := checkURL(svc.host, svc.healthCheck, p.timeout)

	p.lock.Lock()
	defer p.lock.Unlock()

	if err != nil {
		log.Debug().Err(err).Str("host", svc.host).Msg("healtcheck failed")
		if svc.alive {
			log.Warn().Str("host", svc.host).Bool("alive", svc.alive).Msg("changed state")
			svc.alive = false
		}
		return false
	}

	if !svc.alive {
		log.Info().Str("host", svc.host).Bool("alive", svc.alive).Msg("changed state")
		svc.alive = true
	}
	return true
}

func checkURL(baseUrl string, path string, timeout time.Duration) error {
	var resp *http.Response
	var err error

	client := http.Client{Timeout: timeout}

	url := fmt.Sprintf("%s/%s", baseUrl, path)
	resp, err = client.Get(url)

	if err != nil {
		return errors.New(fmt.Sprintf("failed to hit %s with error: %s", url, err.Error()))
	}

	defer func() {
		if e := resp.Body.Close(); e != nil {
			log.Warn().Err(e).Msg("failed to close response body")
		}
	}()

	if resp.StatusCode >= 400 {
		return errors.New(fmt.Sprintf("bad status code %d for %s", resp.StatusCode, url))
	}

	return nil
}
