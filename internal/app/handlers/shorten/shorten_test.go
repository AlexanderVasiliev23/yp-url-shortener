package shorten

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	path         = "/api/shorten"
	defaultToken = "test_token"
	addr         = "localhost:8080"
)

var (
	ErrRepositorySaving = errors.New("repository saving error")
	ErrTokenGen         = errors.New("token gen err")
)

type tokenGeneratorMock struct {
	token string
	err   error
}

func (t tokenGeneratorMock) Generate() (string, error) {
	return t.token, t.err
}

type repositoryMock struct {
	err error
}

func (r repositoryMock) Add(ctx context.Context, shortLink *models.ShortLink) error {
	return r.err
}

func (r repositoryMock) GetTokenByURL(ctx context.Context, url string) (string, error) {
	return defaultToken, nil
}

func TestShorten(t *testing.T) {
	type request struct {
		method string
		body   string
	}

	type want struct {
		code int
		body string
		err  error
	}

	testCases := []struct {
		name           string
		request        request
		want           want
		tokenGenerator tokenGenerator
		repository     repository
	}{
		{
			name: "success",
			request: request{
				method: http.MethodPost,
				body:   `{"url": "https://practicum.yandex.ru/"}`,
			},
			want: want{
				code: http.StatusCreated,
				body: fmt.Sprintf("{\"result\":\"%s/%s\"}\n", addr, defaultToken),
			},
			tokenGenerator: tokenGeneratorMock{token: defaultToken},
			repository:     repositoryMock{err: nil},
		},
		{
			name: "empty url",
			request: request{
				method: http.MethodPost,
				body:   `{"url": ""}`,
			},
			want: want{
				code: http.StatusBadRequest,
				body: "",
				err:  ErrURLIsEmpty,
			},
		},
		{
			name: "wrong request body",
			request: request{
				method: http.MethodPost,
				body:   `{"wrong_field": "v"}`,
			},
			want: want{
				code: http.StatusBadRequest,
				body: "",
				err:  ErrURLIsEmpty,
			},
		},
		{
			name: "token generating error",
			request: request{
				method: http.MethodPost,
				body:   `{"url": "https://practicum.yandex.ru/"}`,
			},
			want: want{
				code: http.StatusInternalServerError,
				body: "",
				err:  ErrTokenGen,
			},
			tokenGenerator: tokenGeneratorMock{err: ErrTokenGen},
		},
		{
			name: "repository saving error",
			request: request{
				method: http.MethodPost,
				body:   `{"url": "https://practicum.yandex.ru/"}`,
			},
			want: want{
				code: http.StatusInternalServerError,
				body: "",
				err:  ErrRepositorySaving,
			},
			tokenGenerator: tokenGeneratorMock{},
			repository:     repositoryMock{err: ErrRepositorySaving},
		},
		{
			name: "already exists",
			request: request{
				method: http.MethodPost,
				body:   `{"url": "https://practicum.yandex.ru/"}`,
			},
			want: want{
				code: http.StatusConflict,
				body: fmt.Sprintf("{\"result\":\"%s/%s\"}\n", addr, defaultToken),
			},
			tokenGenerator: tokenGeneratorMock{},
			repository:     repositoryMock{err: storage.ErrAlreadyExists},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.request.method, path, strings.NewReader(tc.request.body))
			w := httptest.NewRecorder()

			h := Shorten(tc.repository, tc.tokenGenerator, addr)

			e := echo.New()
			c := e.NewContext(r, w)

			err := h(c)

			assert.ErrorIs(t, tc.want.err, err)

			assert.Equal(t, tc.want.code, c.Response().Status)
			assert.Equal(t, tc.want.body, w.Body.String())
		})
	}
}
