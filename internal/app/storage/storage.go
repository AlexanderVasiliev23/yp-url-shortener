package storage

import (
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage/file"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage/local"
	"os"
)

type Storage interface {
	Add(token, url string) error
	Get(token string) (string, error)
}

func New(storageFilePath string) (Storage, error) {
	if storageFilePath == "" {
		return local.New(), nil
	}

	f, err := os.OpenFile(storageFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("can't open storage file: %w", err)
	}

	return file.New(f), nil
}
