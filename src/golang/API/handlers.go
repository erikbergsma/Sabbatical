package main

import (
	"net/http"
	"fmt"
	"strings"
	"strconv"
	"encoding/json"
)

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
	}

	//needs string checking?
	var server Server
	server.Name		= r.FormValue("Name")
	server.Enabled, err	= strconv.ParseBool(r.FormValue("Enabled"))
	server.ID, err		= strconv.Atoi(r.FormValue("ID"))

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

