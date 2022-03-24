package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/pprof"
	"os"
)

func Run() error {
	r := mux.NewRouter()
	r.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	r.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	r.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	r.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	r.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	r.PathPrefix("/debug/").Handler(http.DefaultServeMux)
	// Ports for Heroku
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "ok")
	}))

	return http.ListenAndServe(":"+port, r)
}
