package get

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage/local"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultToken    = "default_test_token"
	defaultSavedURL = "default_saved_url"
	defaultUserID   = 123
)

type mockRepo struct {
	err error
	url *models.ShortLink
}

func (m mockRepo) Get(ctx context.Context, s string) (*models.ShortLink, error) {
	return m.url, m.err
}

func TestGet(t *testing.T) {
	type want struct {
		locationHeader string
		code           int
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
			repo:   mockRepo{url: models.NewShortLink(defaultUserID, uuid.New(), defaultToken, defaultSavedURL)},
			method: http.MethodGet,
			token:  defaultToken,
			want: want{
				code:           http.StatusTemporaryRedirect,
				locationHeader: defaultSavedURL,
			},
		},
		{
			name:   "token not found in repo",
			repo:   mockRepo{err: local.ErrURLNotFound},
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
		{
			name: "deleted url",
			repo: mockRepo{url: func() *models.ShortLink {
				m := models.NewShortLink(defaultUserID, uuid.New(), defaultToken, defaultSavedURL)
				m.Delete()
				return m
			}()},
			method: http.MethodGet,
			token:  defaultToken,
			want:   want{code: http.StatusGone},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(tt.repo).Get

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
