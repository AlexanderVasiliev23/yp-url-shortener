package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url is not found")
)

type localStorage map[string]string

func NewLocalStorage() *localStorage {
	s := make(localStorage)
	return &s
}

func (s localStorage) Add(token, url string) error {
	s[token] = url

	return nil
}

func (s localStorage) Get(token string) (string, error) {
	url, ok := s[token]
	if ok {
		return url, nil
	}

	return "", ErrURLNotFound
}
