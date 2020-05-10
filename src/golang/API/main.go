package main

import (
	"flag"
	"net/http"
	"os"
	log "github.com/sirupsen/logrus"
)

func init(){
	_ = flag.Bool("help", false, "[optional] use: ADDRESS, DB and PASSWORD Environment values to specify a Redis endpoint")
	flag.Parse()

	setupRedisConnection()
	populate()

	log.SetOutput(os.Stdout)
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
