package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/k0kubun/pp"
	"net/http"
	"os"
)

func Run(ctx context.Context) error {
	r := mux.NewRouter()

	//r.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	//r.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	//r.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	//r.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	//r.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	//r.PathPrefix("/debug/").Handler(http.DefaultServeMux)
	// Ports for Heroku
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "ok")
	}))
	server := &http.Server{Addr: ":" + port, Handler: r}
	go func() {
		<-ctx.Done()
		err := server.Shutdown(ctx)
		if err != nil && err != http.ErrServerClosed && err != context.Canceled {
			pp.Println("server.Run: Error:", err)
		}
	}()
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}
	return nil
}
