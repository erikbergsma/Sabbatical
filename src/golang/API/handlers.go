package main

import (
	"net/http"
	"fmt"
	"strings"
	"strconv"
	"encoding/json"
	"io/ioutil"
	"log"
)
type test_struct struct {
    Test string
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
	}

	//needs string checking?
	var server Server
	server.Name		= r.FormValue("Name")
	server.Enabled, err	= strconv.ParseBool(r.FormValue("Enabled"))
	//assign a newly generated/incremented id
	server.ID		= incrGlobalCustomerId()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save to database
	serverToRedis(server)

	http.Redirect(w, r, "/", 301)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("updateHandler")

	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
			panic(err)
	}

	log.Println(string(body))
	var dumpserver ServerString
	//var t test_struct
	err = json.Unmarshal(body, &dumpserver)

	if err != nil {
			panic(err)
	}

	keyname := strings.Join([]string{redisHashKeyRoot, strconv.FormatInt(dumpserver.ID, 10)}, ":")
	fmt.Println(keyname)
	server := getCustomerByKeyname(keyname)

	// needs checking which field needs updating?
	server.Name					= dumpserver.Name
	server.Enabled, err	= strconv.ParseBool(dumpserver.Enabled)

	if err != nil {
		fmt.Println("problem formatting boolean")
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

	ID := strings.TrimSpace(r.FormValue("ID"))
	keyname := strings.Join([]string{redisHashKeyRoot, ID}, ":")
	fmt.Println("deleting: ", keyname)

	// delete the redis hash (customer object)
	err = client.Del(keyname).Err()
	if err != nil {
		fmt.Println("err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// delete key from index
	err = client.SRem(redisSetKeyName, keyname).Err()
	if err != nil {
		fmt.Println("err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", 301)
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
