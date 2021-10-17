package main

import (
	"log"
	"net/http"
)

type statusWriter struct {
	http.ResponseWriter
	statusCode int
}

func (s *statusWriter) WriteHeader(statusCode int) {
	s.ResponseWriter.WriteHeader(statusCode)
	s.statusCode = statusCode
}

type logMiddleware struct {
	http.Handler
}

func (l *logMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s := statusWriter{ResponseWriter: w, statusCode: http.StatusOK}
	defer func() {
		// Log in a defer statement in case the handler panics.
		log.Printf(
			"%s %s -> %d %s",
			req.Method,
			req.URL.String(),
			s.statusCode,
			http.StatusText(s.statusCode),
		)
	}()
	l.Handler.ServeHTTP(&s, req)

}
