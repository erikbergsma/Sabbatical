package main

import (
  "os"
  "fmt"
  log "github.com/sirupsen/logrus"
)

// this sets the loglevel to info if there is nothing in the
// environment or the environment cannot be parsed
func setLogger(){
	lvl, ok := os.LookupEnv("LOGLEVEL")
	// LOG_LEVEL not set, let's default to debug
	if !ok {
			lvl = "info"
	}

	// parse string, this is built-in feature of logrus
	ll, err := log.ParseLevel(lvl)
	if err != nil {
			ll = log.InfoLevel
	}
	// set global log level
	log.SetLevel(ll)

  fmt.Println(log.GetLevel())

	log.SetOutput(os.Stdout)
}
