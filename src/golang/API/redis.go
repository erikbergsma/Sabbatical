package main

import (
	"os"
	"fmt"
	"strings"
	"github.com/fatih/structs"
	"github.com/go-redis/redis/v7"
	"strconv"
	"time"
)

func setupRedisConnection() error {
	addr, ok := os.LookupEnv("ADDRESS")
	if ok != true {
		addr = "localhost:6379"
	}

	fmt.Println("addr:", addr)

	password, ok := os.LookupEnv("PASSWORD")
	if ok != true {
		password = "" // no password set
		fmt.Println("password: <redacted, default>")
	} else {
		fmt.Println("password: <redacted, from ENV>")
	}

	var db int
	dbstring, ok := os.LookupEnv("DB")
	if ok != true {
		db = 0 // use default DB
	} else {
		db, err = strconv.Atoi(dbstring)
		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}
	}
	fmt.Println("db:", db)

	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	retries := 5
	for {
		pong, err := client.Ping().Result()
		fmt.Println(pong, err)

		if err != nil {
			fmt.Println(err)
			retries--

			if retries == 0 {
                                panic(err)
                        }

                        fmt.Println("retrying: ", retries, " more time(s)")
                        time.Sleep(500 * time.Millisecond)

		} else {
			return nil
		}
	}
}

func serverToRedis(server Server) error {
	ID := strconv.FormatInt(server.ID, 10)
	m := structs.Map(server)
	keyname := strings.Join([]string{redisHashKeyRoot, ID}, ":")

	retries := 5
	for {
		success := true
		err := client.HSet(keyname, m).Err()
		err2 := client.SAdd(redisSetKeyName, keyname).Err()

		if err != nil {
			success = false
			fmt.Println("HSET failed: ", err)
		} else if err2 != nil {
			success = false
			fmt.Println("SAdd failed: ", err)
		}

		if success == false {
			retries--

			if retries == 0 {
				return err
			}

			fmt.Println("retrying: ", retries, " more time(s)")
			time.Sleep(500 * time.Millisecond)
		} else {
			return nil
		}
	}
}

func populate() error {
	if len(getCustomersIndex()) != 0 {
		fmt.Println("found entries in DB, skipping the populate process")
		return nil
	}

	id := incrGlobalCustomerId()

	server1 := Server{
		Name:    "gopher3",
		ID:      id,
		Enabled: true,
	}
	serverToRedis(server1)

	id = incrGlobalCustomerId()

	server2 := Server{
		Name:    "gopher6",
		ID:      id,
		Enabled: false,
	}
	serverToRedis(server2)

	return nil
}

func getCustomersIndex() []string{
	allcustomers, err := client.SMembers(redisSetKeyName).Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
	} else {
		fmt.Println("members", allcustomers)
	}

	return allcustomers
}

func getCustomerByKeyname(customer string) Server {
	val, err := client.HGetAll(customer).Result()

	if err == redis.Nil {
		fmt.Println("key does not exist")
	} else if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
	} else {
		fmt.Println("key", val)
	}

	var server Server
	server.Name = val["Name"]

	//this should be gracefull?
	server.Enabled, err = strconv.ParseBool(val["Enabled"])
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("erik was here")
		fmt.Println(err)
		//return
	}

	//this should be gracefull?
	server.ID, err = strconv.ParseInt(val["ID"], 10, 64)
	if err != nil {
		fmt.Println(err)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		//return
	}

	return server
}

func incrGlobalCustomerId() int64 {
	result, err := client.Incr(redisIdKeyName).Result()
	if err != nil {
	    fmt.Println("error incrementing:", err)
	}

	fmt.Println("new id: ", result)
	return result
}
