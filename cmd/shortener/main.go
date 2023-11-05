package main

import (
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/tokenGenerator"
	"net/http"
)

const (
	addr = "localhost:8080"
)

func main() {
	localStorage := storage.NewLocalStorage()
	tokenGenerator := tokenGenerator.New()

	mux := http.NewServeMux()

	handler := handlers.NewHandler(localStorage, tokenGenerator, addr)

	mux.HandleFunc("/", handler.Handle)

	if err := http.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}
}
