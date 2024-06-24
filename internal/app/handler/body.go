package handler

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetBody(r *http.Request) ([]byte, error) {
	fmt.Println(`!strings.Contains(r.Header.Get("Content-Encoding"), "gzip")`, !strings.Contains(r.Header.Get("Content-Encoding"), "gzip"))
	if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		return GetUnCompressedBody(r)
	}

	return GetCompressedBody(r)
}

func GetCompressedBody(r *http.Request) ([]byte, error) {
	gz, err := gzip.NewReader(r.Body)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = gz.Close()
	}()

	// при чтении вернётся распакованный слайс байт
	body, err := io.ReadAll(gz)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func GetUnCompressedBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = r.Body.Close()
	if err != nil {
		return nil, err
	}

	return body, err
}
