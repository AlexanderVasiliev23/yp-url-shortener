package file

import (
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage/storageerrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const (
	defaultToken = "default_token"
	defaultURL   = "default_url"
	filePath     = "/tmp/test_storage.json"
)

func TestFileStorage(t *testing.T) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	require.NoError(t, err)

	storage := New(file)

	url, err := storage.Get(defaultToken)
	assert.Error(t, storageerrors.ErrURLNotFound, err)
	assert.Equal(t, "", url)

	err = storage.Add(defaultToken, defaultURL)
	require.NoError(t, err)

	url, err = storage.Get(defaultToken)
	assert.NoError(t, err)
	assert.Equal(t, defaultURL, url)

	err = os.Remove(filePath)
	require.NoError(t, err)
}
