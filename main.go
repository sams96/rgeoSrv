/*
Command rgeoSrv wraps the package rgeo into a reverse geocoding microservice.

See https://github.com/sams96/rgeo for more information on rgeo.
*/
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sams96/rgeoSrv/query"
)

var (
	srvAddr = ":8080"
)

func main() {
	l := log.New(os.Stderr, "rgeoSrv ", log.LstdFlags)

	mux := http.NewServeMux()
	q, err := query.NewHandlers(l)
	if err != nil {
		l.Fatalf("server failed to start: %v", err)
	}

	q.SetupRoutes(mux)

	srv := newServer(mux, srvAddr)

	l.Println("server starting at", srvAddr)
	err = srv.ListenAndServe()
	if err != nil {
		l.Fatalf("server failed to start: %v", err)
	}
}

// newServer creates a new http server
func newServer(mux *http.ServeMux, addr string) *http.Server {
	srv := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}

	return srv
}
