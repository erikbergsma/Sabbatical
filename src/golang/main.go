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
	server := &Server{
		Name:    "gopher",
		ID:      123456,
		Enabled: true,
	}

	// Convert a struct to a map[string]interface{}
	// => {"Name":"gopher", "ID":123456, "Enabled":true}
	m := structs.Map(server)
	fmt.Println("map", m)

	err := client.HSet("customer:2", m).Err()
	if err != nil {
		panic(err)
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

	val2, err := client.HGetAll("customer:2" ).Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}

	Enabled, err := strconv.ParseBool(val2["Enabled"])
	if err != nil {
		panic(err)
	}
	ID, err := strconv.Atoi(val2["ID"])
	if err != nil {
		panic(err)
	}

	server := &Server{
		Name:    val2["Name"],
		ID:      ID,
		Enabled: Enabled,
	}

	fmt.Println("struct", server)

  w.Header().Set("Server", "A Go Web Server")

	js, err := json.Marshal(server)
	if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

	w.Header().Set("Content-Type", "application/json")
  w.Write(js)

}
