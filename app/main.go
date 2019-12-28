package main

import (
	"github.com/dp0h/srp/app/config"
	"github.com/dp0h/srp/app/pool"
	"github.com/dp0h/srp/app/proxy"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var opts struct {
	Port          int           `long:"port" env:"SRP_PORT" description:"port" default:"443"`
	SslMode       string        `long:"ssl-mode" env:"SRP_SSL_MODE" description:"ssl mode" choice:"none" choice:"static" choice:"auto" default:"none"`
	Host          string        `long:"host" env:"SRP_HOST" description:"host name"`
	CertFile      string        `long:"cert-file" env:"SRP_CERT_FILE" description:"path to cert.pem file"`
	KeyFile       string        `long:"key-file" env:"SRP_KEY_FILE" description:"path to cert.key file"`
	AutoCertPath  string        `long:"autocert-path" env:"SRP_AUTOCERT_PATH" description:"dir where certificates will be stored by autocert manager" default:"./var/autocert"`
	Conf          string        `long:"conf" env:"SRP_CONF" description:"configuration file" default:"srp.yml"`
	Refresh       time.Duration `long:"refresh" env:"SRP_REFRESH" description:"keep alive refresh interval" default:"30s"`
	TimeOut       time.Duration `long:"timeout" env:"SRP_TIMEOUT" description:"keep alive timeouts" default:"10s"`
	ValidateCerts bool          `long:"secure-cert" env:"SRP_VALIDATE_CERTS" description:"validate certificates"`
	Dbg           bool          `long:"dbg" env:"SRP_DEBUG" description:"debug mode"`
}

var revision = "unknown"

func main() {
	log.Printf("SRP - %s", revision)
	setupLog(true)

	if _, err := flags.Parse(&opts); err != nil {
		log.Fatal().Err(err).Msg("failed to parse args")
	}

	confReader, err := os.Open(opts.Conf)
	if err != nil {
		log.Fatal().Err(err).Str("file", opts.Conf).Msg("failed to open config file")
	}

	conf := config.NewConf(confReader)
	if err := confReader.Close(); err != nil {
		log.Warn().Err(err).Str("file", opts.Conf).Msg("failed to close config file")
	}

	rwp := pool.NewRandomWeightedPool(conf, opts.Refresh, opts.TimeOut)
	proxy.NewReverseProxyServer(opts.Port, opts.SslMode, opts.Host, opts.CertFile, opts.KeyFile, opts.AutoCertPath, opts.ValidateCerts, rwp).Run()
}

func setupLog(dbg bool) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if dbg {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
