package main

import (
	"flag"
	"net/http"
	"fmt"
	"github.com/go-redis/redis/v7"
)

var (
	client			= &redis.Client{}
	err			error
	authenticated		= false
	// a list of ids will be stored under customer => [1, 2, 3]
	redisSetKeyName		= "customers"
	// customers will be stored under customer:id => foo
	redisHashKeyRoot	= "customer"
	// customer id's will be generated server side based
	redisIdKeyName          = "lastCustomerId"
)

func main() {
	_ = flag.Bool("help", false, "[optional] use: ADDRESS, DB and PASSWORD Environment values to specify a Redis endpoint")
	flag.Parse()

	setupRedisConnection()
	populate()

	// route
	http.HandleFunc("/list", listHandler)
	//http.HandleFunc("/create", createHandler)
	http.HandleFunc("/update", updateHandler)
	//http.HandleFunc("/delete", deleteHandler)
	fmt.Println("all systems green, launching API  on port 3333")
	http.ListenAndServe(":3333", nil)
}
