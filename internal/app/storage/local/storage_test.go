package local

import (
	"context"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
	"testing"

	"github.com/google/uuid"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/uuidgenerator/google"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/uuidgenerator/mock"

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
	require.NoError(t, err)

	shortLink, err := s.Get(context.Background(), testToken)
	assert.NoError(t, err)
	assert.Equal(t, testURL, shortLink.Original)
}

func TestGetFound(t *testing.T) {
	shortLink := models.NewShortLink(testUserID, uuid.New(), testToken, testURL)

	s := New(google.UUIDGenerator{})
	err := s.Add(context.Background(), shortLink)
	require.NoError(t, err)

	url, err := s.Get(context.Background(), shortLink.Token)

	require.NoError(t, err)
	assert.Equal(t, testURL, url.Original)
}

func TestGetNotFound(t *testing.T) {
	s := New(google.UUIDGenerator{})
	url, err := s.Get(context.Background(), testToken)

	assert.ErrorIs(t, err, ErrURLNotFound)
	assert.Nil(t, url)
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

func TestStats(t *testing.T) {
	s := New(google.UUIDGenerator{})

	shortLink1 := models.NewShortLink(testUserID, uuid.New(), testToken+"1", testURL+"1")
	shortLink2 := models.NewShortLink(testUserID, uuid.New(), testToken+"2", testURL+"2")
	shortLink3 := models.NewShortLink(testUserID+1, uuid.New(), testToken+"3", testURL+"3")

	err := s.Add(context.Background(), shortLink1)
	require.NoError(t, err)

	err = s.Add(context.Background(), shortLink2)
	require.NoError(t, err)

	err = s.Add(context.Background(), shortLink3)
	require.NoError(t, err)

	stats, err := s.Stats(context.Background())
	require.NoError(t, err)

	expected := storage.StatsOutDTO{
		UrlsCount:  3,
		UsersCount: 2,
	}

	require.Equal(t, expected, *stats)
}
