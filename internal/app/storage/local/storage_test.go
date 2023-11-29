package local

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
	s := New()
	err := s.Add(testToken, testURL)

	require.NoError(t, err)
	assert.Equal(t, Storage{testToken: testURL}, *s)
}

func TestGetFound(t *testing.T) {
	s := Storage{testToken: testURL}
	url, err := s.Get(testToken)

	require.NoError(t, err)
	assert.Equal(t, testURL, url)
}

func TestGetNotFound(t *testing.T) {
	s := New()
	url, err := s.Get(testToken)

	assert.ErrorIs(t, err, ErrURLNotFound)
	assert.Equal(t, "", url)
}
