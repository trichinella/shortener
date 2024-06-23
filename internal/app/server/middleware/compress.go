package middleware

import (
	"compress/gzip"
	"go.uber.org/zap"
	"io"
	"net/http"
	"slices"
	"strings"
)

var compressableContentTypes []string

func init() {
	compressableContentTypes = []string{"application/json", "text/html"}
}

func Compress(sugar *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !slices.Contains(compressableContentTypes, r.Header.Get("Content-Type")) {
				next.ServeHTTP(w, r)
				return
			}

			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}

			gz, err := gzip.NewWriterLevel(w, gzip.DefaultCompression)

			if err != nil {
				sugar.Fatal(err)
				panic(err)
			}

			defer func() {
				if err := gz.Close(); err != nil {
					sugar.Fatal(err)
					panic(err)
				}
			}()

			w.Header().Set("Content-Encoding", "gzip")
			next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
		})
	}
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
