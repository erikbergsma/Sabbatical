package main

import (
	"net/http"
	"fmt"
	"html/template"
	"github.com/go-redis/redis"
	"strconv"
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

func main() {
	// route
	http.HandleFunc("/list", listHandler)
	http.ListenAndServe(":8080", nil)
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

	fmt.Println(customers)

	var funcMap = template.FuncMap{
			"multiplication": func(n float64, f float64) float64 {
				return n * f
			},
			"addOne": func(n int) int {
				return n + 1
			},
	}

	t, err := template.New("list.html").Funcs(funcMap).ParseFiles("tmpl/list.html")
	if err != nil {
	  http.Error(w, err.Error(), http.StatusInternalServerError)
	  return
	}

	err = t.Execute(w, customers)
	if err != nil {
	  http.Error(w, err.Error(), http.StatusInternalServerError)
	  return
	}
}
