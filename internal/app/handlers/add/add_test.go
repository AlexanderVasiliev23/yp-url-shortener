package add

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	addr         = "localhost:8080"
	defaultToken = "default_test_token"
)

type mockTokenGenerator struct {
	token string
	err   error
}

func (m mockTokenGenerator) Generate() (string, error) {
	return m.token, m.err
}

type mockRepo struct {
	addingError error
}

func (m mockRepo) Add(_, _ string) error {
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
		tokGen mockTokenGenerator
		method string
		body   string
		want   want
	}{
		{
			name:   "success",
			repo:   mockRepo{},
			tokGen: mockTokenGenerator{token: defaultToken},
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
			tokGen: mockTokenGenerator{},
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
			tokGen: mockTokenGenerator{},
			method: http.MethodPost,
			body:   "test_url",
			want: want{
				code: http.StatusInternalServerError,
				body: "",
			},
		},
		{
			name:   "token generator error",
			repo:   mockRepo{},
			tokGen: mockTokenGenerator{err: errors.New("")},
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
			handler := Add(tt.repo, tt.tokGen, addr)

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
