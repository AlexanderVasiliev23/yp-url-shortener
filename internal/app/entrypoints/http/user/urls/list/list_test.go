package list

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/user/url/list"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
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

type useCaseMock struct {
	outDTO *list.OutDTO
	err    error
}

func (m *useCaseMock) List(ctx context.Context) (*list.OutDTO, error) {
	return m.outDTO, m.err
}

func TestUrls(t *testing.T) {
	type want struct {
		err  error
		body string
		code int
	}

	testCases := []struct {
		name string
		useCase
		want want
	}{
		{
			name: "empty list",
			useCase: &useCaseMock{
				outDTO: nil,
				err:    list.ErrNoSavedURLs,
			},
			want: want{
				code: http.StatusNoContent,
			},
		},
		{
			name: "success list",
			useCase: &useCaseMock{
				outDTO: &list.OutDTO{
					Items: []list.OutDTOItem{
						{
							ShortURL:    fmt.Sprintf("%s/%s", defaultAddr, defaultToken),
							OriginalURL: defaultOriginal,
						},
					},
				},
				err: nil,
			},
			want: want{
				code: http.StatusOK,
				body: fmt.Sprintf(`[{"short_url":"%s/%s","original_url":"%s"}]`+"\n", defaultAddr, defaultToken, defaultOriginal),
			},
		},
		{
			name: "unauthorized",
			useCase: &useCaseMock{
				outDTO: nil,
				err:    list.ErrUnauthorized,
			},
			want: want{
				code: http.StatusUnauthorized,
				err:  list.ErrUnauthorized,
			},
		},
		{
			name: "usecase unknown error",
			useCase: &useCaseMock{
				outDTO: nil,
				err:    ErrDefault,
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

			h := NewHandler(tc.useCase).List

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
