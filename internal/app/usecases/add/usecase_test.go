package add

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/util/auth/mock"
)

const (
	addr         = "localhost:8080"
	defaultToken = "default_test_token"
)

var (
	errDefault = errors.New("test_error")
)

type mockTokenGenerator struct {
	err   error
	token string
}

func (m mockTokenGenerator) Generate() (string, error) {
	return m.token, m.err
}

type mockRepo struct {
	addingErr   error
	getTokenErr error
}

func (m mockRepo) Add(ctx context.Context, shortLink *models.ShortLink) error {
	return m.addingErr
}

func (m mockRepo) GetTokenByURL(ctx context.Context, url string) (string, error) {
	return defaultToken, m.getTokenErr
}

func TestAdd(t *testing.T) {
	type want struct {
		shortURL string
		err      error
	}

	testCases := []struct {
		name               string
		repo               mockRepo
		tokGen             mockTokenGenerator
		userContextFetcher userContextFetcher
		originalURL        string
		want               want
	}{
		{
			name:               "success",
			repo:               mockRepo{},
			tokGen:             mockTokenGenerator{token: defaultToken},
			userContextFetcher: &mock.UserContextFetcherMock{},
			originalURL:        "test_url",
			want: want{
				err:      nil,
				shortURL: fmt.Sprintf("%s/%s", addr, defaultToken),
			},
		},
		{
			name:               "empty originalURL",
			repo:               mockRepo{},
			tokGen:             mockTokenGenerator{},
			userContextFetcher: &mock.UserContextFetcherMock{},
			originalURL:        "",
			want: want{
				err:      ErrOriginalURLIsEmpty,
				shortURL: "",
			},
		},
		{
			name:               "repo returns an error on adding",
			repo:               mockRepo{addingErr: errDefault},
			tokGen:             mockTokenGenerator{},
			userContextFetcher: &mock.UserContextFetcherMock{},
			originalURL:        "test_url",
			want: want{
				err:      errDefault,
				shortURL: "",
			},
		},
		{
			name:               "repo returns an error on getting by token",
			repo:               mockRepo{addingErr: storage.ErrAlreadyExists, getTokenErr: errDefault},
			tokGen:             mockTokenGenerator{},
			userContextFetcher: &mock.UserContextFetcherMock{},
			originalURL:        "test_url",
			want: want{
				err:      errDefault,
				shortURL: "",
			},
		},
		{
			name:               "token generator error",
			repo:               mockRepo{},
			tokGen:             mockTokenGenerator{err: errDefault},
			userContextFetcher: &mock.UserContextFetcherMock{},
			originalURL:        "test_url",
			want: want{
				err:      errDefault,
				shortURL: "",
			},
		},
		{
			name:               "already exists",
			repo:               mockRepo{addingErr: storage.ErrAlreadyExists},
			tokGen:             mockTokenGenerator{},
			userContextFetcher: &mock.UserContextFetcherMock{},
			originalURL:        "test_url",
			want: want{
				err:      ErrOriginURLAlreadyExists,
				shortURL: fmt.Sprintf("%s/%s", addr, defaultToken),
			},
		},
		{
			name:               "fetching userID error",
			repo:               mockRepo{},
			tokGen:             mockTokenGenerator{},
			userContextFetcher: &mock.UserContextFetcherMock{Err: errDefault},
			originalURL:        "test_url",
			want: want{
				err:      errDefault,
				shortURL: "",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			useCase := NewUseCase(tc.repo, tc.tokGen, tc.userContextFetcher, addr)

			shortURL, err := useCase.Add(context.Background(), tc.originalURL)

			if tc.want.err != nil {
				assert.Equal(t, tc.want.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.want.shortURL, shortURL)
		})
	}
}
