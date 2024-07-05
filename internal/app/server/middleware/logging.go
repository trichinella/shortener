package middleware

import (
	"net/http"
	"shortener/internal/app/logging"
	"time"
)

type responseData struct {
	statusCode int
	size       int
}

type LogResponseWriter struct {
	ResponseData *responseData
	http.ResponseWriter
}

func (lrw LogResponseWriter) Write(input []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(input)
	lrw.ResponseData.size = size

	return size, err
}

func (lrw LogResponseWriter) WriteHeader(statusCode int) {
	lrw.ResponseWriter.WriteHeader(statusCode)
	lrw.ResponseData.statusCode = statusCode
}

func LogMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			lrw := LogResponseWriter{
				ResponseData:   &responseData{},
				ResponseWriter: w,
			}
			next.ServeHTTP(lrw, r)
			logging.Sugar.Infow("Request",
				"URI", r.RequestURI,
				"Method", r.Method,
				"Execution time", time.Since(start),
				"Code", lrw.ResponseData.statusCode,
				"Size", lrw.ResponseData.size,
			)
		})
	}
}
