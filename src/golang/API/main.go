package main

import (
	"flag"
	"net/http"
	log "github.com/sirupsen/logrus"
)

func init(){
	_ = flag.Bool("help", false, "[optional] use: ADDRESS, DB and PASSWORD Environment values to specify a Redis endpoint")
	flag.Parse()
	setLogger()
	setupRedisConnection()
	populate()
}

func main() {
	var router = setup_version_1()

	log.Info("all systems green, launching API  on port 33335")
	http.ListenAndServe(":3333", router)
}
