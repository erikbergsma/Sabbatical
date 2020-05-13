package main

import (
	"github.com/go-redis/redis/v7"
)


var (
	client = &redis.Client{}
	err		  error
)

const(
	authenticated			= false
	// a list of ids will be stored under customer => [1, 2, 3]
	redisSetKeyName		= "customers"
	// customers will be stored under customer:id => foo
	redisHashKeyRoot	= "customer"
	// customer id's will be generated server side based
	redisIdKeyName		= "lastCustomerId"
)
