package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var revision = "unknown"

func main() {
	log.Printf("SRP - %s", revision)
	setupLog(true)
}

func setupLog(dbg bool) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if dbg {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
