package main

import (
	"net/http"
	log "github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
)

func setup_version_1() *mux.Router {
  // route
  var router = mux.NewRouter()
  var api = router.PathPrefix("/").Subrouter()

  api.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    w.WriteHeader(http.StatusNotFound)
  })

  api.Use( func(next http.Handler) http.Handler {
    return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request) {
        log.Debug(r.RequestURI)
        next.ServeHTTP(w, r)
    })
  })

  var version1 = api.PathPrefix("/v1").Subrouter()
  add_handlers(version1)

  var latest = api.PathPrefix("/latest").Subrouter()
  add_handlers(latest)

  return router

}

func add_handlers(temp *mux.Router){
  temp.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
      w.WriteHeader(http.StatusOK)
  })

  temp.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      w.WriteHeader(http.StatusForbidden)
  })

  temp.HandleFunc("/list", listHandler)
  temp.HandleFunc("/create", createHandler)
  temp.HandleFunc("/update", updateHandler)
  temp.HandleFunc("/delete", deleteHandler)

}
