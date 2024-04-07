package stats

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
	iputil "github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/util/ip"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	defaultUrlsCount           = 45
	defaultUsersCount          = 23
	defaultTrustedSubnet       = "127.0.0.1/24"
	defaultIPAddress           = "127.0.0.1"
	defaultNotTrustedIPAddress = "123.123.123.123"
)

var (
	errRepo = errors.New("repo error")
)

type repositoryMock struct {
	statusOutDTO *storage.StatsOutDTO
	err          error
}

func (r *repositoryMock) Stats(ctx context.Context) (*storage.StatsOutDTO, error) {
	return r.statusOutDTO, r.err
}

func TestHandle(t *testing.T) {
	type want struct {
		err  error
		body string
		code int
	}

	testCases := []struct {
		name string
		ip   string
		repo *repositoryMock
		want want
	}{
		{
			name: "success",
			ip:   defaultIPAddress,
			repo: &repositoryMock{
				statusOutDTO: &storage.StatsOutDTO{
					UrlsCount:  defaultUrlsCount,
					UsersCount: defaultUsersCount,
				},
				err: nil,
			},
			want: want{
				err:  nil,
				body: fmt.Sprintf(`{"urls":%d,"users":%d}`, defaultUrlsCount, defaultUsersCount) + "\n",
				code: http.StatusOK,
			},
		},
		{
			name: "invalid ip address",
			ip:   "invalid ip address",
			want: want{
				err:  fmt.Errorf("invalid ip address: %s", "invalid ip address"),
				body: "",
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "not trusted ip address",
			ip:   defaultNotTrustedIPAddress,
			want: want{
				err:  nil,
				body: "",
				code: http.StatusForbidden,
			},
		},
		{
			name: "repo error",
			ip:   defaultIPAddress,
			repo: &repositoryMock{
				err: errRepo,
			},
			want: want{
				err:  errRepo,
				body: "",
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			request.Header.Set(iputil.IPHeader, tc.ip)

			h := NewHandler(tc.repo, defaultTrustedSubnet).Handle

			e := echo.New()
			c := e.NewContext(request, recorder)

			err := h(c)

			if tc.want.err == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.want.err.Error())
			}

			assert.Equal(t, tc.want.code, recorder.Code)
			assert.Equal(t, tc.want.body, recorder.Body.String())
		})
	}
}
