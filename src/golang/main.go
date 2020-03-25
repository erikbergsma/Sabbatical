package main

import (
	"net/http"
	"fmt"
	"strings"
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
	authenticated 		= false
	redisSetKeyName		= "customers"
	redisHashKeyRoot	=	"customer"
)

type Server struct {
	Name        string
	ID          int
	Enabled     bool
	users       []string // not exported
}

func serverToRedis(server Server) error {
	m := structs.Map(server)
	ID := strconv.Itoa(m["ID"].(int))
	keyname := strings.Join([]string{redisHashKeyRoot, ID}, ":")

  err := client.HSet(keyname, m).Err()
	err2 := client.SAdd(redisSetKeyName, keyname).Err()

	if err != nil {
		return err
	} else if err2 != nil {
		return err2
	}

	return nil
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

func main() {
	populate()

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

	allcustomers, err := client.SMembers(redisSetKeyName).Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("members", allcustomers)
	}

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

		//this should be gracefull?
		Enabled, err := strconv.ParseBool(val2["Enabled"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
	    return
		}

		//this should be gracefull?
		ID, err := strconv.Atoi(val2["ID"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
	    return
		}

		// i need to somehow verify if the types are correct
		// or check for errors when creating this here type
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
