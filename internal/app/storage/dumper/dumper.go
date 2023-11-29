package dumper

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
	"github.com/google/uuid"
	"io"
	"os"
)

type record struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (r record) isValid() bool {
	return r.UUID != "" && r.ShortURL != "" && r.OriginalURL != ""
}

type Storage struct {
	wrappedStorage storage.Storage
	file           *os.File
	notSyncedYet   []*record
	bufferSize     int
}

func New(wrappedStorage storage.Storage, filepath string, bufferSize int) (*Storage, error) {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("opening storage file: %w", err)
	}

	s := &Storage{
		wrappedStorage: wrappedStorage,
		file:           file,
		notSyncedYet:   []*record{},
		bufferSize:     bufferSize,
	}

	if err := s.recoverDataFromFile(); err != nil {
		return nil, fmt.Errorf("recovering storage data from file %w", err)
	}

	return s, nil
}

func (s *Storage) Add(token, url string) error {
	if err := s.wrappedStorage.Add(token, url); err != nil {
		return fmt.Errorf("adding to wrapped storage: %w", err)
	}

	record := record{
		UUID:        uuid.NewString(),
		ShortURL:    token,
		OriginalURL: url,
	}

	s.notSyncedYet = append(s.notSyncedYet, &record)

	if len(s.notSyncedYet) > s.bufferSize {
		if err := s.Dump(); err != nil {
			return fmt.Errorf("dump records on adding: %w", err)
		}
	}

	return nil
}

func (s *Storage) Get(token string) (string, error) {
	return s.wrappedStorage.Get(token)
}

func (s *Storage) Dump() error {
	encoder := json.NewEncoder(s.file)

	for _, r := range s.notSyncedYet {
		if err := encoder.Encode(r); err != nil {
			return fmt.Errorf("record incoding: %w", err)
		}
	}

	s.notSyncedYet = []*record{}

	return nil
}

func (s *Storage) recoverDataFromFile() error {
	_, err := s.file.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("storage file seek: %w", err)
	}

	scanner := bufio.NewScanner(s.file)

	for scanner.Scan() {
		record := new(record)
		if err := json.Unmarshal(scanner.Bytes(), record); err != nil {
			return fmt.Errorf("unmarshal record: %w", err)
		}

		if !record.isValid() {
			return fmt.Errorf("unmarshalled record is not valid, original row: %s", scanner.Text())
		}

		if err := s.wrappedStorage.Add(record.ShortURL, record.OriginalURL); err != nil {
			return fmt.Errorf("adding to wrapped storage: %w", err)
		}
	}

	return nil
}
