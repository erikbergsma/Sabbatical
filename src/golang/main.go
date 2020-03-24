package main

import (
	"net/http"
	"fmt"
	"github.com/fatih/structs"
	"github.com/go-redis/redis"
	//"reflect"
	//"github.com/mitchellh/mapstructure"
	"strconv"
	"encoding/json"
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

type Server struct {
	Name        string
	ID          int
	Enabled     bool
	users       []string // not exported
}

func main() {
	server1 := &Server{
		Name:    "gopher",
		ID:      123456,
		Enabled: true,
	}

	server2 := &Server{
		Name:    "blooper",
		ID:      78910,
		Enabled: false,
	}

	// Convert a struct to a map[string]interface{}
	// => {"Name":"gopher", "ID":123456, "Enabled":true}
	m := structs.Map(server1)
	fmt.Println("map", m)

	o := structs.Map(server2)
	fmt.Println("map", o)

	err := client.HSet("customer:2", m).Err()
	if err != nil {
		panic(err)
	}

	err2 := client.SAdd("customers", "customer:2").Err()
	if err2 != nil {
		panic(err2)
	}

	err3 := client.HSet("customer:3", o).Err()
	if err3 != nil {
		panic(err3)
	}

	err4 := client.SAdd("customers", "customer:3").Err()
	if err4 != nil {
		panic(err4)
	}

	// route
	http.HandleFunc("/list", listHandler)
	//http.HandleFunc("/create", createHandler)
	//http.HandleFunc("/update", updateHandler)
	//http.HandleFunc("/delete", deleteHandler)
	http.ListenAndServe(":3333", nil)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}

	allcustomers, err := client.SMembers("customers").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("members", allcustomers)
	}

	//customers := make([]struct{})
	var customers []Server
	for _, customer := range allcustomers {
		val2, err := client.HGetAll(customer).Result()
		if err == redis.Nil {
			fmt.Println("key2 does not exist")
		} else if err != nil {
			panic(err)
		} else {
			fmt.Println("key2", val2)
		}

		Enabled, err := strconv.ParseBool(val2["Enabled"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
	    return
		}

		ID, err := strconv.Atoi(val2["ID"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
	    return
		}

		server := Server{
			Name:    val2["Name"],
			ID:      ID,
			Enabled: Enabled,
		}

		customers = append(customers, server)
	}

	w.Header().Set("Server", "A Go Web Server")

	js, err := json.Marshal(customers)
	if err != nil {
	  http.Error(w, err.Error(), http.StatusInternalServerError)
	  return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
