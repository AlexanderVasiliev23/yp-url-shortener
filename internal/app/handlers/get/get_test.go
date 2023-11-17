package get

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultToken    = "default_test_token"
	defaultSavedURL = "default_saved_url"
)

type mockRepo struct {
	err error
	url string
}

func (m mockRepo) Get(s string) (url string, err error) {
	return m.url, m.err
}

func TestGet(t *testing.T) {
	type want struct {
		code           int
		locationHeader string
	}

	tests := []struct {
		name   string
		repo   mockRepo
		method string
		token  string
		want   want
	}{
		{
			name:   "success",
			repo:   mockRepo{url: defaultSavedURL},
			method: http.MethodGet,
			token:  defaultToken,
			want: want{
				code:           http.StatusTemporaryRedirect,
				locationHeader: defaultSavedURL,
			},
		},
		{
			name:   "token not found in repo",
			repo:   mockRepo{err: storage.ErrURLNotFound},
			method: http.MethodGet,
			token:  defaultToken,
			want:   want{code: http.StatusBadRequest},
		},
		{
			name:   "empty token",
			repo:   mockRepo{},
			method: http.MethodGet,
			token:  "",
			want:   want{code: http.StatusBadRequest},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := Get(tt.repo)

			r := httptest.NewRequest(tt.method, "/", nil)
			w := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(r, w)
			c.SetPath("/:token")
			c.SetParamNames("token")
			c.SetParamValues(tt.token)
			err := handler(c)

			require.NoError(t, err)
			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.locationHeader, w.Header().Get("Location"))
		})
	}
}
