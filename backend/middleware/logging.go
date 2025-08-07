package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logging middleware logs HTTP requests with timing information
func Logging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		responseWriter := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next(responseWriter, r)

		// Log the request details
		duration := time.Since(start)
		log.Printf(
			"%s %s %s - %d - %v",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			responseWriter.statusCode,
			duration,
		)
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggingWrapper wraps a handler with logging
func LoggingWrapper(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Logging(handler.ServeHTTP)(w, r)
	})
}

