/*
rgeoSrv wraps the package rgeo into a reverse geocoding microservice.

see https://github.com/sams96/rgeo for more information on rgeo.
*/
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sams96/rgeo"
)

var (
	srvAddr = ":8080"
)

func main() {
	l := log.New(os.Stderr, "rgeoSrv ", log.LstdFlags)
	mux := http.NewServeMux()
	mux.HandleFunc("/", logger(rootHandler, l))

	srv := newServer(mux, srvAddr)

	l.Println("server starting at", srvAddr)
	err := srv.ListenAndServe()
	if err != nil {
		l.Fatalf("server failed to start: %v", err)
	}
}

// rootHandler handles request to "/"
func rootHandler(w http.ResponseWriter, r *http.Request) {
	coords, err := parseURL(r.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loc, err := rgeo.ReverseGeocode(coords)
	if err != nil { //&& err.Error() != "country not found" {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	/* Use with next version of rgeo
	switch err {
	case rgeo.ErrCountryNotFound:
		// Don't return an internal server error, maybe some other error
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	*/

	resp, err := json.Marshal(loc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// parseURL extracts the coordinates from the request URL
func parseURL(u *url.URL) ([]float64, error) {
	if u.RawQuery == "" || !strings.ContainsRune(u.RawQuery, '&') {
		return []float64{}, errors.New("rgeoSrv requires a request in the form /?lon&lat")
	}

	coord := strings.Split(u.RawQuery, "&")

	lon, err := strconv.ParseFloat(coord[0], 64)
	if err != nil {
		return []float64{}, err
	}

	lat, err := strconv.ParseFloat(coord[1], 64)
	if err != nil {
		return []float64{}, err
	}

	return []float64{lon, lat}, nil
}

// logger is http handler middleware which adds logging
func logger(next http.HandlerFunc, l *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next(w, r)
		l.Printf("Request for %s from %s processed in %s\n",
			r.URL.Path, r.RemoteAddr, time.Since(startTime))
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
