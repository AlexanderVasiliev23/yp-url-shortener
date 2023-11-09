package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testToken = "test_token"
	testURL   = "test_url"
)

func TestAdd(t *testing.T) {
	storage := NewLocalStorage()
	err := storage.Add(testToken, testURL)

	require.NoError(t, err)
	assert.Equal(t, LocalStorage{testToken: testURL}, *storage)
}

func TestGetFound(t *testing.T) {
	storage := LocalStorage{testToken: testURL}
	url, err := storage.Get(testToken)

	require.NoError(t, err)
	assert.Equal(t, testURL, url)
}

func TestGetNotFound(t *testing.T) {
	storage := NewLocalStorage()
	url, err := storage.Get(testToken)

	assert.ErrorIs(t, err, ErrURLNotFound)
	assert.Equal(t, "", url)
}
