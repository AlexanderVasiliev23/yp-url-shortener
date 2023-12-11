package dumper

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage/local"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const (
	testStorageFilePath = "/tmp/testStorageFilePath.json"
	defaultBufferSize   = 3
)

type mockStorage struct {
	data map[string]string
}

func (m mockStorage) Add(ctx context.Context, token, url string) error {
	m.data[token] = url
	return nil
}

func (m mockStorage) Get(ctx context.Context, token string) (string, error) {
	url, ok := m.data[token]
	if !ok {
		return "", local.ErrURLNotFound
	}

	return url, nil
}

func (s mockStorage) SaveBatch(ctx context.Context, shortLinks []*models.ShortLink) error {
	return nil
}

func TestStorage_RecoveringFromFileSuccess(t *testing.T) {
	defer os.Remove(testStorageFilePath)

	token := "mbQTUSzkAa"
	URL := "https://ya.ru"

	err := os.WriteFile(testStorageFilePath, []byte(fmt.Sprintf(`{"id":"%s","token":"%s","original":"%s"}`, uuid.NewString(), token, URL)+"\n"), os.ModePerm)
	require.NoError(t, err)

	s, err := New(context.Background(), mockStorage{make(map[string]string)}, testStorageFilePath, defaultBufferSize)
	require.NoError(t, err)

	actualURL, err := s.Get(context.Background(), token)
	require.NoError(t, err)
	assert.Equal(t, URL, actualURL)
}

func TestStorage_RecoveringFromFileFailed(t *testing.T) {
	defer os.Remove(testStorageFilePath)

	err := os.WriteFile(testStorageFilePath, []byte(`{"wrong":"data"}`+"\n"), os.ModePerm)
	require.NoError(t, err)

	_, err = New(context.Background(), mockStorage{make(map[string]string)}, testStorageFilePath, defaultBufferSize)
	require.Error(t, err)
}

func TestStorage_BufferSizeEqualsZero(t *testing.T) {
	defer os.Remove(testStorageFilePath)

	s, err := New(context.Background(), mockStorage{make(map[string]string)}, testStorageFilePath, 0)
	require.NoError(t, err)

	require.NoError(t, s.Add(context.Background(), "token1", "url1"))

	value, err := os.ReadFile(testStorageFilePath)
	require.NoError(t, err)

	assert.Equal(t, 1, rowsInContent(value))
}

func TestStorage_BufferSizeEqualsOne(t *testing.T) {
	defer os.Remove(testStorageFilePath)

	s, err := New(context.Background(), mockStorage{make(map[string]string)}, testStorageFilePath, 1)
	require.NoError(t, err)

	require.NoError(t, s.Add(context.Background(), "token1", "url1"))
	value, err := os.ReadFile(testStorageFilePath)
	require.NoError(t, err)

	assert.Equal(t, 0, rowsInContent(value))

	require.NoError(t, s.Add(context.Background(), "token2", "url2"))
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
