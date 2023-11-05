package handlers

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	addr            = "localhost:8080"
	defaultToken    = "default_test_token"
	defaultSavedURL = "default_saved_url"
)

type mockRepo struct {
	addingError  error
	gettingError error
	url          string
}

func (m mockRepo) Add(token, url string) error {
	return m.addingError
}

func (m mockRepo) Get(s string) (url string, err error) {
	return m.url, m.gettingError
}

type mockTokenGenerator struct {
}

func (m mockTokenGenerator) Generate() string {
	return defaultToken
}

func Test_handler_add(t *testing.T) {
	type want struct {
		code int
		body string
	}

	tests := []struct {
		name   string
		repo   mockRepo
		method string
		body   string
		want   want
	}{
		{
			name:   "success",
			repo:   mockRepo{},
			method: http.MethodPost,
			body:   "test_url",
			want: want{
				code: http.StatusCreated,
				body: fmt.Sprintf("http://%s/%s", addr, defaultToken),
			},
		},
		{
			name:   "empty body",
			repo:   mockRepo{},
			method: http.MethodPost,
			body:   "",
			want: want{
				code: http.StatusBadRequest,
				body: "",
			},
		},
		{
			name:   "repo returns an error",
			repo:   mockRepo{addingError: errors.New("")},
			method: http.MethodPost,
			body:   "test_url",
			want: want{
				code: http.StatusInternalServerError,
				body: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenGenerator := mockTokenGenerator{}

			handler := NewHandler(tt.repo, tokenGenerator, addr)

			r := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			handler.Handle(w, r)

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.body, w.Body.String())
		})
	}
}

func Test_handler_get(t *testing.T) {
	type want struct {
		code           int
		locationHeader string
	}

	tests := []struct {
		name   string
		repo   mockRepo
		method string
		path   string
		want   want
	}{
		{
			name:   "success",
			repo:   mockRepo{url: defaultSavedURL},
			method: http.MethodGet,
			path:   fmt.Sprintf("/%s", defaultToken),
			want: want{
				code:           http.StatusTemporaryRedirect,
				locationHeader: defaultSavedURL,
			},
		},
		{
			name:   "token not found in repo",
			repo:   mockRepo{},
			method: http.MethodGet,
			path:   "/",
			want:   want{code: http.StatusBadRequest},
		},
		{
			name:   "empty token",
			repo:   mockRepo{},
			method: http.MethodGet,
			path:   "/",
			want:   want{code: http.StatusBadRequest},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(tt.repo, mockTokenGenerator{}, addr)

			r := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			handler.Handle(w, r)

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.locationHeader, w.Header().Get("Location"))
		})
	}
}
