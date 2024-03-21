package add

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/util/auth/mock"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		body string
		code int
	}

	tests := []struct {
		name               string
		repo               mockRepo
		tokGen             mockTokenGenerator
		userContextFetcher userContextFetcher
		method             string
		body               string
		want               want
	}{
		{
			name:               "success",
			repo:               mockRepo{},
			tokGen:             mockTokenGenerator{token: defaultToken},
			userContextFetcher: &mock.UserContextFetcherMock{},
			method:             http.MethodPost,
			body:               "test_url",
			want: want{
				code: http.StatusCreated,
				body: fmt.Sprintf("%s/%s", addr, defaultToken),
			},
		},
		{
			name:               "empty body",
			repo:               mockRepo{},
			tokGen:             mockTokenGenerator{},
			userContextFetcher: &mock.UserContextFetcherMock{},
			method:             http.MethodPost,
			body:               "",
			want: want{
				code: http.StatusBadRequest,
				body: "",
			},
		},
		{
			name:               "repo returns an error on adding",
			repo:               mockRepo{addingErr: errDefault},
			tokGen:             mockTokenGenerator{},
			userContextFetcher: &mock.UserContextFetcherMock{},
			method:             http.MethodPost,
			body:               "test_url",
			want: want{
				code: http.StatusInternalServerError,
				body: "",
			},
		},
		{
			name:               "repo returns an error on getting by token",
			repo:               mockRepo{addingErr: storage.ErrAlreadyExists, getTokenErr: errDefault},
			tokGen:             mockTokenGenerator{},
			userContextFetcher: &mock.UserContextFetcherMock{},
			method:             http.MethodPost,
			body:               "test_url",
			want: want{
				code: http.StatusInternalServerError,
				body: "",
			},
		},
		{
			name:               "token generator error",
			repo:               mockRepo{},
			tokGen:             mockTokenGenerator{err: errDefault},
			userContextFetcher: &mock.UserContextFetcherMock{},
			method:             http.MethodPost,
			body:               "test_url",
			want: want{
				code: http.StatusInternalServerError,
				body: "",
			},
		},
		{
			name:               "already exists",
			repo:               mockRepo{addingErr: storage.ErrAlreadyExists},
			tokGen:             mockTokenGenerator{},
			userContextFetcher: &mock.UserContextFetcherMock{},
			method:             http.MethodPost,
			body:               "test_url",
			want: want{
				code: http.StatusConflict,
				body: fmt.Sprintf("%s/%s", addr, defaultToken),
			},
		},
		{
			name:               "fetching userID error",
			repo:               mockRepo{},
			tokGen:             mockTokenGenerator{},
			userContextFetcher: &mock.UserContextFetcherMock{Err: errDefault},
			method:             http.MethodPost,
			body:               "test_url",
			want: want{
				code: http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(tt.repo, tt.tokGen, tt.userContextFetcher, addr).Add

			r := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(r, w)

			err := handler(c)

			if tt.want.code == http.StatusCreated {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.body, w.Body.String())
		})
	}
}
