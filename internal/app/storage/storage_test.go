package storage

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	testToken = "test_token"
	testURL   = "test_url"
)

func TestAdd(t *testing.T) {
	storage := NewLocalStorage()
	err := storage.Add(testToken, testURL)

	require.NoError(t, err)
	assert.Equal(t, localStorage{testToken: testURL}, *storage)
}

func TestGetFound(t *testing.T) {
	storage := localStorage{testToken: testURL}
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
