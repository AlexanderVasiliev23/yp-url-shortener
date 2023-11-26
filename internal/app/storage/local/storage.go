package local

import (
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage/storage_errors"
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

	return "", storage_errors.ErrURLNotFound
}
