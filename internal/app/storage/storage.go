package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url is not found")
)

type LocalStorage map[string]string

func NewLocalStorage() *LocalStorage {
	s := make(LocalStorage)
	return &s
}

func (s LocalStorage) Add(token, url string) error {
	s[token] = url

	return nil
}

func (s LocalStorage) Get(token string) (string, error) {
	url, ok := s[token]
	if ok {
		return url, nil
	}

	return "", ErrURLNotFound
}
