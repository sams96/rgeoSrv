package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.RawQuery == "" || !strings.ContainsRune(r.URL.RawQuery, '&') {
		http.Error(w, "rgeoSrv requires a request in the form /?lon&lat",
			http.StatusInternalServerError)
		return
	}

	coord := strings.Split(r.URL.RawQuery, "&")

	lon, err := strconv.ParseFloat(coord[0], 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lat, err := strconv.ParseFloat(coord[1], 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loc, err := rgeo.ReverseGeocode([]float64{lon, lat})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

func logger(next http.HandlerFunc, l *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next(w, r)
		l.Printf("Request for %s from %s processed in %s\n",
			r.URL.Path, r.RemoteAddr, time.Since(startTime))
	}
}

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
