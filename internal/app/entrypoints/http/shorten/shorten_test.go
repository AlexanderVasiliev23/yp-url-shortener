package shorten

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/shorten/single"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
)

const (
	path         = "/api/shorten"
	defaultToken = "test_token"
	addr         = "https://my_url_shortener"
)

var (
	errDefault = errors.New("default error")
)

type useCaseMock struct {
	shortURL string
	err      error
}

func (m *useCaseMock) Shorten(ctx context.Context, jsonString string) (shortURL string, err error) {
	return m.shortURL, m.err
}

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
	type request struct {
		method string
		body   string
	}

	type want struct {
		err  error
		body string
		code int
	}

	testCases := []struct {
		useCase useCase
		request request
		name    string
		want    want
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
			useCase: &useCaseMock{
				shortURL: fmt.Sprintf("%s/%s", addr, defaultToken),
				err:      nil,
			},
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
				err:  single.ErrEmptyOriginalURL,
			},
			useCase: &useCaseMock{
				shortURL: "",
				err:      single.ErrEmptyOriginalURL,
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
				err:  single.ErrEmptyOriginalURL,
			},
			useCase: &useCaseMock{
				shortURL: "",
				err:      single.ErrEmptyOriginalURL,
			},
		},
		{
			name: "usecase unknown error",
			request: request{
				method: http.MethodPost,
				body:   `{"url": "https://practicum.yandex.ru/"}`,
			},
			want: want{
				code: http.StatusInternalServerError,
				body: "",
				err:  errDefault,
			},
			useCase: &useCaseMock{
				shortURL: "",
				err:      errDefault,
			},
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
				err:  single.ErrAlreadyExists,
			},
			useCase: &useCaseMock{
				shortURL: fmt.Sprintf("%s/%s", addr, defaultToken),
				err:      single.ErrAlreadyExists,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.request.method, path, strings.NewReader(tc.request.body))
			w := httptest.NewRecorder()

			h := NewShortener(tc.useCase).Handle

			e := echo.New()
			c := e.NewContext(r, w)

			err := h(c)

			assert.ErrorIs(t, tc.want.err, err)

			assert.Equal(t, tc.want.code, c.Response().Status)
			assert.Equal(t, tc.want.body, w.Body.String())
		})
	}
}
