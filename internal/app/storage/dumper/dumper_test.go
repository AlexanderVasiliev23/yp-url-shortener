package dumper

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage/local"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/uuidgenerator/google"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const (
	testStorageFilePath = "/tmp/testStorageFilePath.json"
	defaultBufferSize   = 3
	defaultUserID       = 123
)

type mockStorage struct {
	data map[string]*models.ShortLink
}

func (m mockStorage) Add(ctx context.Context, shortLink *models.ShortLink) error {
	m.data[shortLink.Token] = shortLink
	return nil
}

func (m mockStorage) Get(ctx context.Context, token string) (*models.ShortLink, error) {
	url, ok := m.data[token]
	if !ok {
		return nil, local.ErrURLNotFound
	}

	return url, nil
}

func (m mockStorage) SaveBatch(ctx context.Context, shortLinks []*models.ShortLink) error {
	return nil
}

func (m mockStorage) GetTokenByURL(ctx context.Context, url string) (string, error) {
	return "", nil
}

func (m mockStorage) FindByUserID(ctx context.Context, userID int) ([]*models.ShortLink, error) {
	return nil, nil
}

func (m mockStorage) DeleteTokens(ctx context.Context, userID int, tokens []string) error {
	return nil
}

func TestStorage_RecoveringFromFileSuccess(t *testing.T) {
	defer os.Remove(testStorageFilePath)

	token := "mbQTUSzkAa"
	URL := "https://ya.ru"

	err := os.WriteFile(testStorageFilePath, []byte(fmt.Sprintf(`{"id":"%s","token":"%s","original":"%s","user_id":%d}`, uuid.NewString(), token, URL, defaultUserID)+"\n"), os.ModePerm)
	require.NoError(t, err)

	s, err := New(context.Background(), mockStorage{make(map[string]*models.ShortLink)}, google.UUIDGenerator{}, testStorageFilePath, defaultBufferSize)
	require.NoError(t, err)

	actualURL, err := s.Get(context.Background(), token)
	require.NoError(t, err)
	assert.Equal(t, URL, actualURL.Original)
}

func TestStorage_RecoveringFromFileFailed(t *testing.T) {
	defer os.Remove(testStorageFilePath)

	err := os.WriteFile(testStorageFilePath, []byte(`{"wrong":"data"}`+"\n"), os.ModePerm)
	require.NoError(t, err)

	_, err = New(context.Background(), mockStorage{make(map[string]*models.ShortLink)}, google.UUIDGenerator{}, testStorageFilePath, defaultBufferSize)
	require.Error(t, err)
}

func TestStorage_BufferSizeEqualsZero(t *testing.T) {
	defer os.Remove(testStorageFilePath)

	s, err := New(context.Background(), mockStorage{make(map[string]*models.ShortLink)}, google.UUIDGenerator{}, testStorageFilePath, 0)
	require.NoError(t, err)

	require.NoError(t, s.Add(context.Background(), models.NewShortLink(defaultUserID, uuid.New(), "token1", "url1")))

	value, err := os.ReadFile(testStorageFilePath)
	require.NoError(t, err)

	assert.Equal(t, 1, rowsInContent(value))
}

func TestStorage_BufferSizeEqualsOne(t *testing.T) {
	defer os.Remove(testStorageFilePath)

	s, err := New(context.Background(), mockStorage{make(map[string]*models.ShortLink)}, google.UUIDGenerator{}, testStorageFilePath, 1)
	require.NoError(t, err)

	require.NoError(t, s.Add(context.Background(), models.NewShortLink(defaultUserID, uuid.New(), "token1", "url1")))
	value, err := os.ReadFile(testStorageFilePath)
	require.NoError(t, err)

	assert.Equal(t, 0, rowsInContent(value))

	require.NoError(t, s.Add(context.Background(), models.NewShortLink(defaultUserID, uuid.New(), "token2", "url2")))
	value, err = os.ReadFile(testStorageFilePath)
	require.NoError(t, err)

	assert.Equal(t, 2, rowsInContent(value))
}

func rowsInContent(content []byte) int {
	result := 0
	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		result++
	}

	return result
}
