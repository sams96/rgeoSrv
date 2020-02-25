/*
Copyright 2020 Sam Smith

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License.  You may obtain a copy of the
License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied.  See the License for the
specific language governing permissions and limitations under the License.
*/

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
will return:
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
