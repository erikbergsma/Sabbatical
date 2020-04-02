package main

import (
	"net/http"
	"fmt"
	"strings"
	"github.com/fatih/structs"
	"github.com/go-redis/redis/v7"
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
	ID := strconv.Itoa(server.ID)
	m := structs.Map(server)
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
	http.HandleFunc("/create", createHandler)
	//http.HandleFunc("/update", updateHandler)
	//http.HandleFunc("/delete", deleteHandler)
	fmt.Println("all systems green, launching API  on port 3333")
	http.ListenAndServe(":3333", nil)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
	}

	//needs string checking?
	var server Server
	server.Name					= r.FormValue("Name")
	server.Enabled, err	= strconv.ParseBool(r.FormValue("Enabled"))
	server.ID, err 			= strconv.Atoi(r.FormValue("ID"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save to database
	serverToRedis(server)

	http.Redirect(w, r, "/", 301)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
	}

	keyname := strings.Join([]string{redisHashKeyRoot, r.FormValue("ID")}, ":")
	server := getCustomerByKeyname(keyname)

	// needs checking which field needs updating?
	server.Name					= r.FormValue("Name")
	server.Enabled, err	= strconv.ParseBool(r.FormValue("Enabled"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save to database
	serverToRedis(server)

	http.Redirect(w, r, "/", 301)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
	}

	keyname := strings.Join([]string{redisHashKeyRoot, r.FormValue("ID")}, ":")

  err = client.Del(keyname).Err()
	if err != nil {
		fmt.Println("err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = client.SRem(redisSetKeyName, keyname).Err()
	if err != nil {
		fmt.Println("err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", 301)
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

func listHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}

	allcustomers := getCustomersIndex()

	var customers []Server
	for _, customer := range allcustomers {
		server := getCustomerByKeyname(customer)
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
