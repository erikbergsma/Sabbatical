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
	// route
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/delete", deleteHandler)

	log.Info("all systems green, launching API  on port 3333")
	http.ListenAndServe(":3333", nil)
}
