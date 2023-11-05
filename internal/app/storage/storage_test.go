package storage

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	testToken = "test_token"
	testUrl   = "test_url"
)

func TestAdd(t *testing.T) {
	storage := NewLocalStorage()
	err := storage.Add(testToken, testUrl)

	require.NoError(t, err)
	assert.Equal(t, localStorage{testToken: testUrl}, *storage)
}

func TestGetFound(t *testing.T) {
	storage := localStorage{testToken: testUrl}
	url, err := storage.Get(testToken)

	require.NoError(t, err)
	assert.Equal(t, testUrl, url)
}

func TestGetNotFound(t *testing.T) {
	storage := NewLocalStorage()
	url, err := storage.Get(testToken)

	assert.ErrorIs(t, err, ErrURLNotFound)
	assert.Equal(t, "", url)
}
