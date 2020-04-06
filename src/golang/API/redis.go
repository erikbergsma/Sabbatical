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
	ID := strconv.Itoa(server.ID)
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

func populate(){
	server1 := Server{
		Name:    "gopher3",
		ID:      111458,
		Enabled: true,
	}
	serverToRedis(server1)

	server2 := Server{
		Name:    "gopher6",
		ID:      1232213,
		Enabled: false,
	}
	serverToRedis(server2)
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
	val2, err := client.HGetAll(customer).Result()

	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
	} else {
		fmt.Println("key2", val2)
	}

	var server Server
	server.Name = val2["Name"]

	//this should be gracefull?
	server.Enabled, err = strconv.ParseBool(val2["Enabled"])
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		//return
	}

	//this should be gracefull?
	server.ID, err = strconv.Atoi(val2["ID"])
	if err != nil {
		fmt.Println(err)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		//return
	}

	return server
}
