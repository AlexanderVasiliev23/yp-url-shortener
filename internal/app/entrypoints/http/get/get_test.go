package get

import (
	"context"
	"errors"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/get"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const (
	defaultToken    = "default_test_token"
	defaultSavedURL = "default_saved_url"
	defaultUserID   = 123
)

var (
	errDefault = errors.New("default error")
)

type useCaseMock struct {
	originalURL string
	err         error
}

func (m *useCaseMock) Get(ctx context.Context, token string) (originalURL string, err error) {
	return m.originalURL, m.err
}

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
		err            error
	}

	testCases := []struct {
		name    string
		method  string
		token   string
		useCase useCase
		want    want
	}{
		{
			name:   "success",
			method: http.MethodGet,
			token:  defaultToken,
			useCase: &useCaseMock{
				originalURL: defaultSavedURL,
				err:         nil,
			},
			want: want{
				code:           http.StatusTemporaryRedirect,
				locationHeader: defaultSavedURL,
				err:            nil,
			},
		},
		{
			name:   "useCase unknown error",
			method: http.MethodGet,
			token:  defaultToken,
			useCase: &useCaseMock{
				originalURL: "",
				err:         errDefault,
			},
			want: want{
				code: http.StatusInternalServerError,
				err:  errDefault,
			},
		},
		{
			name:   "empty token",
			method: http.MethodGet,
			token:  "",
			useCase: &useCaseMock{
				originalURL: "",
				err:         get.ErrTokenIsEmpty,
			},
			want: want{
				code: http.StatusBadRequest,
				err:  get.ErrTokenIsEmpty,
			},
		},
		{
			name:   "deleted url",
			method: http.MethodGet,
			token:  defaultToken,
			useCase: &useCaseMock{
				originalURL: "",
				err:         get.ErrTokenIsDeleted,
			},
			want: want{
				code: http.StatusGone,
				err:  get.ErrTokenIsDeleted,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewHandler(tc.useCase).Get

			r := httptest.NewRequest(tc.method, "/", nil)
			w := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(r, w)
			c.SetPath("/:token")
			c.SetParamNames("token")
			c.SetParamValues(tc.token)
			err := handler(c)

			if tc.want.err != nil {
				assert.Equal(t, tc.want.err, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.want.code, w.Code)
			assert.Equal(t, tc.want.locationHeader, w.Header().Get("Location"))
		})
	}
}
