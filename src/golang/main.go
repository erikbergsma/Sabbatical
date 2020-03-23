package main

import (
	"net/http"
	"fmt"
	"github.com/fatih/structs"
	"github.com/go-redis/redis"
)

var (
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err error
	authenticated = false
)


func main() {
	err := client.Set("key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get("key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	type Server struct {
		Name        string
		ID          int
		Enabled     bool
		users       []string // not exported
	}

	server := &Server{
		Name:    "gopher",
		ID:      123456,
		Enabled: true,
	}

	// Convert a struct to a map[string]interface{}
	// => {"Name":"gopher", "ID":123456, "Enabled":true}
	m := structs.Map(server)
	fmt.Println("map", m)

	client.HSet("customer:2", m)

	val3 := client.HGetAll("customer:2" )
	fmt.Println("key3", val3)

	val4 := client.HGet("customer:2", "Name")
	fmt.Println("key4", val4)

	// route
	//http.HandleFunc("/list", listHandler)
	//http.HandleFunc("/create", createHandler)
	//http.HandleFunc("/update", updateHandler)
	//http.HandleFunc("/delete", deleteHandler)
	//http.ListenAndServe(":3333", nil)

}


func listHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}

	val2, err := client.Get("key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}

}
