package main

import (
	"net/http"
	"strings"
	"strconv"
	"encoding/json"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
)
type test_struct struct {
    Test string
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var newcustomer Customer
	err = json.Unmarshal(body, &newcustomer)

	if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//we create the ID serverside
	newcustomer.ID		= incrGlobalCustomerId()

	// Save to database
	customerToRedis(newcustomer)

	//return to the client so he can fetch the ID / check
	js, err := json.Marshal(newcustomer)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var dumpcustomer Customer
	err = json.Unmarshal(body, &dumpcustomer)

	if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	keyname := strings.Join([]string{redisHashKeyRoot, strconv.FormatInt(dumpcustomer.ID, 10)}, ":")
	customer := getCustomerByKeyname(keyname)

	// needs checking which field needs updating?
	customer.Name					= dumpcustomer.Name
  customer.Enabled	    = dumpcustomer.Enabled

	// Save to database
	customerToRedis(customer)

	http.Redirect(w, r, "/", 301)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var dumpcustomer Customer
	err = json.Unmarshal(body, &dumpcustomer)

	if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	keyname := strings.Join([]string{redisHashKeyRoot, strconv.FormatInt(dumpcustomer.ID, 10)}, ":")
	log.Debug("deleting: ", keyname)

	// delete the redis hash (customer object)
	err = client.Del(keyname).Err()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// delete key from index
	err = client.SRem(redisSetKeyName, keyname).Err()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", 301)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	//with the new gorilla/mux router, this does not work
	//and /or the Method is set to HEAD..... ?

	//if r.Method != "GET" {
	//	http.Error(w, "Method not allowed", http.StatusBadRequest)
	//}

	allcustomers := getCustomersIndex()

	var customers []Customer
	for _, customer := range allcustomers {
		customer := getCustomerByKeyname(customer)
		customers = append(customers, customer)
	}

	w.Header().Set("Customer", "A Go Web Customer")

	js, err := json.Marshal(customers)
	if err != nil {
		log.Error(err.Error())
	  http.Error(w, err.Error(), http.StatusInternalServerError)
	  return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
