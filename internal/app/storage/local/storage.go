package local

import (
	"errors"
)

var (
	ErrURLNotFound = errors.New("url is not found")
)

type Storage map[string]string

func New() *Storage {
	s := make(Storage)
	return &s
}

func (s Storage) Add(token, url string) error {
	s[token] = url

	return nil
}

func (s Storage) Get(token string) (string, error) {
	url, ok := s[token]
	if ok {
		return url, nil
	}

	return "", ErrURLNotFound
}
