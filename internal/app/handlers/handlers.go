package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type repository interface {
	Add(token, url string) error
	Get(string) (url string, err error)
}

type tokenGenerator interface {
	Generate() string
}

type handler struct {
	repository     repository
	tokenGenerator tokenGenerator
	addr           string
}

func NewHandler(repository repository, tokenGenerator tokenGenerator, addr string) *handler {
	return &handler{repository: repository, tokenGenerator: tokenGenerator, addr: addr}
}

func (h handler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.get(w, r)
	case http.MethodPost:
		h.add(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h handler) add(w http.ResponseWriter, r *http.Request) {
	url, err := io.ReadAll(r.Body)
	if err != nil || len(url) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token := h.tokenGenerator.Generate()
	if err := h.repository.Add(token, string(url)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = fmt.Fprintf(w, "http://%s/%s", h.addr, token)
}

func (h handler) get(w http.ResponseWriter, r *http.Request) {
	token := strings.Trim(r.URL.Path, "/")

	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url, err := h.repository.Get(token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
