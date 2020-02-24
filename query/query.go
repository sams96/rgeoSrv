/*
Package query implements the query handler for rgeoSrv
*/
package query

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/sams96/rgeo"
)

// Handlers supplies dependencies for query handlers
type Handlers struct {
	l *log.Logger
	r *rgeo.Rgeo
}

// NewHandlers creates new Handlers type for query
func NewHandlers(l *log.Logger) (*Handlers, error) {
	r, err := rgeo.New(rgeo.Countries10)
	if err != nil {
		return nil, err
	}

	return &Handlers{
		l: l,
		r: r,
	}, nil
}

// query handles request to "/query"
func (h *Handlers) query(w http.ResponseWriter, r *http.Request) {
	coords, err := parseURL(r.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loc, err := h.r.ReverseGeocode(coords)
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

// parseURL extracts the coordinates from the request URL
func parseURL(u *url.URL) ([]float64, error) {
	if u.RawQuery == "" || !strings.ContainsRune(u.RawQuery, '&') {
		return []float64{},
			errors.New("rgeoSrv requires a request in the form /query?lon&lat")
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
func (h *Handlers) logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next(w, r)
		h.l.Printf("Request for %s from %s processed in %s\n",
			r.URL.Path, r.RemoteAddr, time.Since(startTime))
	}
}

// SetupRoutes sets up the query routes on the given mux
func (h *Handlers) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/query", h.logger(h.query))
}
