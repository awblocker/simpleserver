package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

var (
	port = flag.Int("port", 8000, "Port for server")
	dir  = flag.String("dir", ".", "Directory to serve")
)

func logRequest(t time.Time, r *http.Request) {
	// Handling username with "-" default
	username := "-"
	if r.URL.User != nil {
		if n := r.URL.User.Username(); n != "" {
			username = n
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}

	log.Printf("%s - %s [%s] \"%s %s %s\" - -",
		host,
		username,
		t.Format("02/Jan/2006:15:04:05 -0700"),
		r.Method,
		r.URL.RequestURI(),
		r.Proto,
	)
}

func Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		logRequest(t, r)
		h.ServeHTTP(w, r)
	})
}

func main() {
	flag.Parse()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port),
		Log(http.FileServer(http.Dir(*dir)))))
}
