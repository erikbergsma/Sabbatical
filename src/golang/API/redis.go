package main

import (
	"os"
	"strings"
	"github.com/fatih/structs"
	"github.com/go-redis/redis/v7"
	"strconv"
	"time"
	log "github.com/sirupsen/logrus"
)

func setupRedisConnection() error {
	addr, ok := os.LookupEnv("ADDRESS")
	if ok != true {
		addr = "localhost:6379"
	}

	log.Debug("addr:", addr)

	password, ok := os.LookupEnv("PASSWORD")
	if ok != true {
		password = "" // no password set
		log.Debug("password: <redacted, default>")
	} else {
		log.Debug("password: <redacted, from ENV>")
	}

	var db int
	dbstring, ok := os.LookupEnv("DB")
	if ok != true {
		db = 0 // use default DB
	} else {
		db, err = strconv.Atoi(dbstring)
		if err != nil {
			// handle error
			log.Error(err)
			os.Exit(2)
		}
	}
	log.Debug("db:", db)

	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	retries := 5
	for {
		pong, err := client.Ping().Result()
		log.Debug(pong, err)

		if err != nil {
			log.Error(err)
			retries--

			if retries == 0 {
				panic(err)
			}

			log.Debug("retrying: ", retries, " more time(s)")
			time.Sleep(500 * time.Millisecond)

		} else {
			return nil
		}
	}
}

func customerToRedis(customer Customer) error {
	ID := strconv.FormatInt(customer.ID, 10)
	m := structs.Map(customer)
	keyname := strings.Join([]string{redisHashKeyRoot, ID}, ":")

	retries := 5
	for {
		success := true
		err := client.HSet(keyname, m).Err()
		err2 := client.SAdd(redisSetKeyName, keyname).Err()

		if err != nil {
			success = false
			log.Error("HSET failed: ", err)
		} else if err2 != nil {
			success = false
			log.Error("SAdd failed: ", err)
		}

		if success == false {
			retries--

			if retries == 0 {
				return err
			}

			log.Debug("retrying: ", retries, " more time(s)")
			time.Sleep(500 * time.Millisecond)
		} else {
			return nil
		}
	}
}

func populate() error {
	if len(getCustomersIndex()) != 0 {
		log.Debug("found entries in DB, skipping the populate process")
		return nil
	}

	id := incrGlobalCustomerId()

	customer1 := Customer{
		Name:    "gopher3",
		ID:      id,
		Enabled: true,
	}
	customerToRedis(customer1)

	id = incrGlobalCustomerId()

	customer2 := Customer{
		Name:    "gopher6",
		ID:      id,
		Enabled: false,
	}
	customerToRedis(customer2)

	return nil
}

func getCustomersIndex() []string{
	allcustomers, err := client.SMembers(redisSetKeyName).Result()
	if err == redis.Nil {
		log.Warning("No members in redis list: ", redisSetKeyName)
	} else if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
	} else {
		log.Debug("members", allcustomers)
	}

	return allcustomers
}

func getCustomerByKeyname(customerName string) Customer {
	val, err := client.HGetAll(customerName).Result()

	if err == redis.Nil {
		log.Warning("key does not exist")
	} else if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
	} else {
		log.Debug("key", val)
	}

	var customer Customer
	customer.Name = val["Name"]

	//this should be gracefull?
	customer.Enabled, err = strconv.ParseBool(val["Enabled"])
	if err != nil {
		log.Error(err)
	}

	//this should be gracefull?
	customer.ID, err = strconv.ParseInt(val["ID"], 10, 64)
	if err != nil {
		log.Error(err)
	}

	return customer
}

func incrGlobalCustomerId() int64 {
	result, err := client.Incr(redisIdKeyName).Result()
	if err != nil {
	    log.Error("error incrementing:", err)
	}

	log.Debug("new id: ", result)
	return result
}
