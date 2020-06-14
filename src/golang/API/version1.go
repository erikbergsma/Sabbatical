package main

import (
	"net/http"
	log "github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
)

func setup_version_1() *mux.Router {
  // route
  var router = mux.NewRouter()
  var api = router.PathPrefix("/api").Subrouter()

  api.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    w.WriteHeader(http.StatusNotFound)
  })

  api.Use( func(next http.Handler) http.Handler {
    return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request) {
        log.Debug(r.RequestURI)
        next.ServeHTTP(w, r)
    })
  })

  var api1 = api.PathPrefix("/v1").Subrouter()

  api1.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
      w.WriteHeader(http.StatusOK)
  })

  api1.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      w.WriteHeader(http.StatusForbidden)
  })

  api1.HandleFunc("/list", listHandler)
  api1.HandleFunc("/create", createHandler)
  api1.HandleFunc("/update", updateHandler)
  api1.HandleFunc("/delete", deleteHandler)

  return router

}
