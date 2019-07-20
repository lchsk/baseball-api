package main

import (
	"log"
	"net/http"
	"strings"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		w.Header().Add("Content-Type", "application/json")

		sw := statusWriter{ResponseWriter: w}
		start := time.Now()
		next.ServeHTTP(&sw, r)
		duration := time.Now().Sub(start)

		log.Println(sw.status, r.Method, r.RequestURI, duration.Seconds()*1000)
	})
}
