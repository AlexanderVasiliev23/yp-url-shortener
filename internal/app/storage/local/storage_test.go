package local

import (
	"context"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/uuidgenerator/google"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/uuidgenerator/mock"
	"github.com/google/uuid"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testToken  = "test_token"
	testURL    = "test_url"
	testUserID = 123
)

func TestAdd(t *testing.T) {
	s := New(google.UUIDGenerator{})

	expectedModel := models.NewShortLink(testUserID, uuid.New(), testToken, testURL)

	err := s.Add(context.Background(), expectedModel)

	url, err := s.Get(context.Background(), testToken)
	require.NoError(t, err)

	assert.Equal(t, testURL, url)
}

func TestGetFound(t *testing.T) {
	shortLink := models.NewShortLink(testUserID, uuid.New(), testToken, testURL)

	s := New(google.UUIDGenerator{})
	err := s.Add(context.Background(), shortLink)
	require.NoError(t, err)

	url, err := s.Get(context.Background(), shortLink.Token)

	require.NoError(t, err)
	assert.Equal(t, testURL, url)
}

func TestGetNotFound(t *testing.T) {
	s := New(google.UUIDGenerator{})
	url, err := s.Get(context.Background(), testToken)

	assert.ErrorIs(t, err, ErrURLNotFound)
	assert.Equal(t, "", url)
}

func TestFindByUserID(t *testing.T) {
	resUUID := uuid.New()
	s := New(mock.NewGenerator(resUUID))

	shortLink := models.NewShortLink(testUserID, resUUID, testToken, testURL)
	err := s.Add(context.Background(), shortLink)
	require.NoError(t, err)

	anotherUserID := 234
	shortLinkByAnotherUser := models.NewShortLink(anotherUserID, resUUID, "anotherToken", "anotherURL")
	err = s.Add(context.Background(), shortLinkByAnotherUser)
	require.NoError(t, err)

	actual, err := s.FindByUserID(context.Background(), testUserID)
	assert.NoError(t, err)

	assert.Equal(t, []*models.ShortLink{shortLink}, actual)
}
