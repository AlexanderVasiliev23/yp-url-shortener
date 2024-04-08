package list

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
)

const (
	defaultToken    = "test_token"
	defaultAddr     = "https://my_url_shortener"
	defaultUserID   = 1234
	defaultOriginal = "test_original"
)

var (
	ErrDefault = errors.New("test_error")
)

type userContextFetcherMock struct {
	err    error
	userID int
}

func (f userContextFetcherMock) GetUserIDFromContext(ctx context.Context) (int, error) {
	return f.userID, f.err
}

type storageMock struct {
	err    error
	result []*models.ShortLink
}

func (s storageMock) FindByUserID(ctx context.Context, userID int) ([]*models.ShortLink, error) {
	return s.result, s.err
}

func TestUrls(t *testing.T) {
	type want struct {
		err  error
		body string
		code int
	}

	testCases := []struct {
		userContextFetcherMock *userContextFetcherMock
		name                   string
		storage                storageMock
		want                   want
	}{
		{
			name: "empty list",
			userContextFetcherMock: &userContextFetcherMock{
				userID: defaultUserID,
			},
			storage: storageMock{
				result: make([]*models.ShortLink, 0),
			},
			want: want{
				code: http.StatusNoContent,
			},
		},
		{
			name: "success list",
			userContextFetcherMock: &userContextFetcherMock{
				userID: defaultUserID,
			},
			storage: storageMock{
				result: []*models.ShortLink{
					{
						Token:    defaultToken,
						Original: defaultOriginal,
					},
				},
			},
			want: want{
				code: http.StatusOK,
				body: fmt.Sprintf(`[{"short_url":"%s/%s","original_url":"%s"}]`+"\n", defaultAddr, defaultToken, defaultOriginal),
			},
		},
		{
			name: "storage error",
			userContextFetcherMock: &userContextFetcherMock{
				userID: defaultUserID,
			},
			storage: storageMock{
				err: ErrDefault,
			},
			want: want{
				code: http.StatusInternalServerError,
				err:  ErrDefault,
			},
		},
		{
			name: "unauthorized",
			userContextFetcherMock: &userContextFetcherMock{
				err: ErrDefault,
			},
			want: want{
				code: http.StatusUnauthorized,
				err:  ErrDefault,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/", nil)

			h := NewHandler(tc.storage, tc.userContextFetcherMock, defaultAddr).List

			e := echo.New()
			c := e.NewContext(request, recorder)

			err := h(c)

			if tc.want.err == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, tc.want.err, err)
			}

			assert.Equal(t, tc.want.code, recorder.Code)
			assert.Equal(t, tc.want.body, recorder.Body.String())
		})
	}
}
