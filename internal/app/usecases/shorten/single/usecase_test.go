package single

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/util/auth/mock"
)

const (
	path         = "/api/shorten"
	defaultToken = "test_token"
	addr         = "https://my_url_shortener"
)

var (
	ErrRepositorySaving = errors.New("repository saving error")
	ErrTokenGen         = errors.New("token gen err")
	errDefault          = errors.New("default error")
)

type tokenGeneratorMock struct {
	err   error
	token string
}

func (t tokenGeneratorMock) Generate() (string, error) {
	return t.token, t.err
}

type repositoryMock struct {
	addingErr  error
	gettingErr error
}

func (r repositoryMock) Add(ctx context.Context, shortLink *models.ShortLink) error {
	return r.addingErr
}

func (r repositoryMock) GetTokenByURL(ctx context.Context, url string) (string, error) {
	return defaultToken, r.gettingErr
}

func TestShorten(t *testing.T) {
	type want struct {
		err      error
		shortURL string
	}

	testCases := []struct {
		tokenGenerator     tokenGenerator
		repository         repository
		userContextFetcher userContextFetcher
		originalURL        string
		name               string
		want               want
	}{
		{
			name:        "success",
			originalURL: "https://practicum.yandex.ru",
			want: want{
				err:      nil,
				shortURL: fmt.Sprintf("%s/%s", addr, defaultToken),
			},
			tokenGenerator:     tokenGeneratorMock{token: defaultToken},
			repository:         repositoryMock{},
			userContextFetcher: &mock.UserContextFetcherMock{},
		},
		{
			name:        "empty url",
			originalURL: "",
			want: want{
				shortURL: "",
				err:      ErrEmptyOriginalURL,
			},
			userContextFetcher: &mock.UserContextFetcherMock{},
		},
		{
			name:        "token generating error",
			originalURL: "https://practicum.yandex.ru",
			want: want{
				shortURL: "",
				err:      ErrTokenGen,
			},
			tokenGenerator:     tokenGeneratorMock{err: ErrTokenGen},
			userContextFetcher: &mock.UserContextFetcherMock{},
		},
		{
			name:        "repository saving error",
			originalURL: "https://practicum.yandex.ru",
			want: want{
				shortURL: "",
				err:      ErrRepositorySaving,
			},
			tokenGenerator:     tokenGeneratorMock{},
			repository:         repositoryMock{addingErr: ErrRepositorySaving},
			userContextFetcher: &mock.UserContextFetcherMock{},
		},
		{
			name:        "repository getting error",
			originalURL: "https://practicum.yandex.ru",
			want: want{
				shortURL: "",
				err:      errDefault,
			},
			tokenGenerator:     tokenGeneratorMock{},
			repository:         repositoryMock{addingErr: storage.ErrAlreadyExists, gettingErr: errDefault},
			userContextFetcher: &mock.UserContextFetcherMock{},
		},
		{
			name:        "already exists",
			originalURL: "https://practicum.yandex.ru",
			want: want{
				shortURL: fmt.Sprintf("%s/%s", addr, defaultToken),
				err:      ErrAlreadyExists,
			},
			tokenGenerator:     tokenGeneratorMock{},
			repository:         repositoryMock{addingErr: storage.ErrAlreadyExists},
			userContextFetcher: &mock.UserContextFetcherMock{},
		},
		{
			name:        "user fetcher error",
			originalURL: "https://practicum.yandex.ru",
			want: want{
				shortURL: "",
				err:      errDefault,
			},
			tokenGenerator:     tokenGeneratorMock{},
			repository:         repositoryMock{},
			userContextFetcher: &mock.UserContextFetcherMock{Err: errDefault},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			useCase := NewUseCase(tc.repository, tc.tokenGenerator, tc.userContextFetcher, addr)

			originalURL, err := useCase.Shorten(context.Background(), tc.originalURL)

			if tc.want.err != nil {
				assert.ErrorIs(t, err, tc.want.err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.want.shortURL, originalURL)
		})
	}
}
