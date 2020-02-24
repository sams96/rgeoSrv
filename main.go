/*
Command rgeoSrv wraps the package rgeo into a reverse geocoding microservice.

See https://github.com/sams96/rgeo for more information on rgeo.

Installation

	go get github.com/sams96/rgeoSrv/..
or,
	docker pull docker.pkg.github.com/sams96/rgeosrv/rgeosrv

Usage

	rgeoSrv -addr localhost:8080
or,
	docker run -p 8080:8080 docker.pkg.github.com/sams96/rgeosrv/rgeosrv
and then:
	curl "localhost:8080/query?0&52"
will yield:
	{"country":"United Kingdom","country_long":"United Kingdom of Great Britain and Northern Ireland","country_code_2":"GB","country_code_3":"GBR","continent":"Europe","region":"Europe","subregion":"Northern Europe"}
*/
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sams96/rgeoSrv/query"
)

var (
	defaultSrvAddr = ":8080"
)

func main() {
	// Bind flags
	srvAddr := flag.String("addr", defaultSrvAddr, "Address to open the server on")
	flag.Parse()

	// Create logger
	l := log.New(os.Stderr, "rgeoSrv ", log.LstdFlags)

	// Create mux and add handlers
	mux := http.NewServeMux()
	q, err := query.NewHandlers(l)
	if err != nil {
		l.Fatalf("server failed to start: %v", err)
	}

	q.SetupRoutes(mux)

	// Create and start server
	srv := newServer(mux, *srvAddr)

	l.Println("server starting at", *srvAddr)
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
