package main

import (
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var opts struct {
	Port         int    `long:"port" env:"SRP_PORT" description:"port" default:"443"`
	SslMode      string `long:"ssl-mode" env:"SRP_SSL_MODE" description:"ssl mode" choice:"none" choice:"static" choice:"auto" default:"none"`
	Site         string `long:"site" env:"SRP_SITE" description:"site name"`
	CertFile     string `long:"cert-file" env:"SRP_CERT_FILE" description:"path to cert.pem file"`
	KeyFile      string `long:"key-file" env:"SRP_KEY_FILE" description:"path to cert.key file"`
	AutoCertPath string `long:"autocert-path" env:"SRP_AUTOCERT_PATH" description:"dir where certificates will be stored by autocert manager" default:"./var/autocert"`
	Dbg          bool   `long:"dbg" env:"SRP_DEBUG" description:"debug mode"`
}

var revision = "unknown"

func main() {
	log.Printf("SRP - %s", revision)
	setupLog(true)

	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}


	//

}

func setupLog(dbg bool) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if dbg {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
