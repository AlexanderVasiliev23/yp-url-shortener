package list

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	defaultToken    = "test_token"
	defaultAddr     = "test_addr"
	defaultUserID   = 1234
	defaultOriginal = "test_original"
)

var (
	ErrDefault = errors.New("test_error")
)

type userContextFetcherMock struct {
	userID int
	err    error
}

func (f userContextFetcherMock) GetUserIDFromContext(ctx context.Context) (int, error) {
	return f.userID, f.err
}

type storageMock struct {
	result []*models.ShortLink
	err    error
}

func (s storageMock) FindByUserID(ctx context.Context, userID int) ([]*models.ShortLink, error) {
	return s.result, s.err
}

func TestUrls(t *testing.T) {

	type want struct {
		code int
		body string
		err  error
	}

	testCases := []struct {
		name                   string
		storage                storageMock
		userContextFetcherMock *userContextFetcherMock
		want                   want
	}{
		{
			name: "empty list",
			storage: storageMock{
				result: make([]*models.ShortLink, 0),
			},
			want: want{
				code: http.StatusNoContent,
			},
		},
		{
			name: "success list",
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
			storage: storageMock{
				err: ErrDefault,
			},
			want: want{
				code: http.StatusInternalServerError,
				err:  ErrDefault,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/", nil)

			h := List(tc.storage, userContextFetcherMock{userID: defaultUserID}, defaultAddr)

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