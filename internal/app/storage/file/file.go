package file

import (
	"bufio"
	"encoding/json"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage/storageerrors"
	"github.com/google/uuid"
	"io"
	"os"
	"strings"
)

type Storage struct {
	file *os.File
}

type record struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func New(file *os.File) *Storage {
	return &Storage{file: file}
}

func (s Storage) Add(token, url string) error {
	record := record{
		UUID:        uuid.NewString(),
		ShortURL:    token,
		OriginalURL: url,
	}

	return json.NewEncoder(s.file).Encode(record)
}

func (s Storage) Get(token string) (string, error) {
	_, err := s.file.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(s.file)

	for scanner.Scan() {
		if !strings.Contains(scanner.Text(), token) {
			continue
		}

		record := new(record)
		if err := json.Unmarshal(scanner.Bytes(), record); err != nil {
			return "", err
		}

		return record.OriginalURL, nil
	}

	return "", storageerrors.ErrURLNotFound
}
