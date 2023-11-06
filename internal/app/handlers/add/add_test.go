package add

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	addr         = "localhost:8080"
	defaultToken = "default_test_token"
)

type mockTokenGenerator struct {
}

func (m mockTokenGenerator) Generate() string {
	return defaultToken
}

type mockRepo struct {
	addingError error
	url         string
}

func (m mockRepo) Add(token, url string) error {
	return m.addingError
}

func TestAdd(t *testing.T) {
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
				body: fmt.Sprintf("%s/%s", addr, defaultToken),
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

			handler := Add(tt.repo, tokenGenerator, addr)

			r := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(r, w)

			err := handler(c)

			require.NoError(t, err)
			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.body, w.Body.String())
		})
	}
}
