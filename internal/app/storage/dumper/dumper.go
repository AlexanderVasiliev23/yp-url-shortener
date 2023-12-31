package dumper

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
	"io"
	"os"
)

var _ storage.Storage = (*Storage)(nil)

type Storage struct {
	wrappedStorage storage.Storage
	file           *os.File
	notSyncedYet   []*models.ShortLink
	bufferSize     int
}

func New(ctx context.Context, wrappedStorage storage.Storage, filepath string, bufferSize int) (*Storage, error) {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("opening storage file: %w", err)
	}

	s := &Storage{
		wrappedStorage: wrappedStorage,
		file:           file,
		notSyncedYet:   []*models.ShortLink{},
		bufferSize:     bufferSize,
	}

	if err := s.recoverDataFromFile(ctx); err != nil {
		return nil, fmt.Errorf("recovering storage data from file %w", err)
	}

	return s, nil
}

func (s *Storage) Add(ctx context.Context, token, url string) error {
	if err := s.wrappedStorage.Add(ctx, token, url); err != nil {
		return fmt.Errorf("adding to wrapped storage: %w", err)
	}

	shortLink := models.NewShortLink(token, url)

	s.notSyncedYet = append(s.notSyncedYet, shortLink)

	if len(s.notSyncedYet) > s.bufferSize {
		if err := s.Dump(); err != nil {
			return fmt.Errorf("dump records on adding: %w", err)
		}
	}

	return nil
}

func (s *Storage) Get(ctx context.Context, token string) (string, error) {
	return s.wrappedStorage.Get(ctx, token)
}

func (s *Storage) SaveBatch(ctx context.Context, shortLinks []*models.ShortLink) error {
	for _, shortLink := range shortLinks {
		if err := s.Add(ctx, shortLink.Token, shortLink.Original); err != nil {
			return fmt.Errorf("add one short link: %w", err)
		}
	}

	return nil
}

func (s *Storage) GetTokenByURL(ctx context.Context, url string) (string, error) {
	return s.wrappedStorage.GetTokenByURL(ctx, url)
}

func (s *Storage) Dump() error {
	encoder := json.NewEncoder(s.file)

	for _, r := range s.notSyncedYet {
		if err := encoder.Encode(r); err != nil {
			return fmt.Errorf("record incoding: %w", err)
		}
	}

	s.notSyncedYet = []*models.ShortLink{}

	return nil
}

func (s *Storage) recoverDataFromFile(ctx context.Context) error {
	_, err := s.file.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("storage file seek: %w", err)
	}

	scanner := bufio.NewScanner(s.file)

	for scanner.Scan() {
		shortLink := new(models.ShortLink)
		if err := json.Unmarshal(scanner.Bytes(), shortLink); err != nil {
			return fmt.Errorf("unmarshal record: %w", err)
		}

		if !shortLink.IsValid() {
			return fmt.Errorf("unmarshalled record is not valid, original row: %s", scanner.Text())
		}

		if err := s.wrappedStorage.Add(ctx, shortLink.Token, shortLink.Original); err != nil {
			return fmt.Errorf("adding to wrapped storage: %w", err)
		}
	}

	return nil
}
