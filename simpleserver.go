package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	port = flag.Int("port", 8000, "Port for server")
	dir  = flag.String("dir", ".", "Directory to serve")
)

type LoggingResponseWriter struct {
	writer http.ResponseWriter
	status int
	size   int
}

func (lw *LoggingResponseWriter) Write(b []byte) (int, error) {
	n, err := lw.writer.Write(b)
	lw.size += n
	return n, err
}

func (lw *LoggingResponseWriter) Header() http.Header {
	return lw.writer.Header()
}

func (lw *LoggingResponseWriter) WriteHeader(i int) {
	lw.status = i
	lw.writer.WriteHeader(i)
}

func logHTTP(t time.Time, lw *LoggingResponseWriter, r *http.Request) {
	// Handling username with "-" default
	username := "-"
	if r.URL.User != nil {
		if n := r.URL.User.Username(); n != "" {
			username = n
		}
	}

	log.Printf("%s - %s [%s] \"%s %s %s\" %d %d",
		r.Host,
		username,
		t.Format("02/Jan/2006:15:04:05 -0700"),
		r.Method,
		r.URL.RequestURI(),
		r.Proto,
		lw.status,
		lw.size,
	)
}

func Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lw := &LoggingResponseWriter{writer: w}
		t := time.Now()
		defer logHTTP(t, lw, r)
		h.ServeHTTP(lw, r)
	})
}

func main() {
	flag.Parse()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port),
		Log(http.FileServer(http.Dir(*dir)))))
}
