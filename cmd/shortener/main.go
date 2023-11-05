package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	addr = "localhost:8080"
)

var (
	ErrURLNotFound = errors.New("url is not found")
)

type localStorage map[string]string

func newStorage() localStorage {
	return make(localStorage)
}

func (s localStorage) add(token, url string) {
	s[token] = url
}

func (s localStorage) get(token string) (string, error) {
	url, ok := s[token]
	if ok {
		return url, nil
	}

	return "", ErrURLNotFound
}

func main() {
	storage := newStorage()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			token := strings.Trim(r.URL.Path, "/")

			url, err := storage.get(token)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			w.Header().Set("Location", url)
			w.WriteHeader(http.StatusTemporaryRedirect)

			return
		}

		if r.Method == http.MethodPost {
			url, err := io.ReadAll(r.Body)
			if err != nil || len(url) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			token := "EwHXdJfB"
			storage.add(token, string(url))

			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, "http://%s/%s", addr, token)

			return
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	if err := http.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}
}
