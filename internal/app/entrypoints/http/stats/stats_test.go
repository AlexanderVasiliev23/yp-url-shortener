package stats

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/stats"
	iputil "github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/util/ip"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	defaultUrlsCount  = 45
	defaultUsersCount = 23
)

var (
	errDefault = errors.New("default error")
)

type useCaseMock struct {
	outDTO *stats.OutDTO
	err    error
}

func (m *useCaseMock) Stats(ctx context.Context, ip string) (*stats.OutDTO, error) {
	return m.outDTO, m.err
}

func TestHandle(t *testing.T) {
	type want struct {
		err  error
		body string
		code int
	}

	testCases := []struct {
		name    string
		ip      string
		useCase useCase
		want    want
	}{
		{
			name: "success",
			useCase: &useCaseMock{
				outDTO: &stats.OutDTO{
					Urls:  defaultUrlsCount,
					Users: defaultUsersCount,
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
			name: "not trusted ip address",
			useCase: &useCaseMock{
				outDTO: nil,
				err:    stats.ErrNotTrustedIP,
			},
			want: want{
				err:  nil,
				body: "",
				code: http.StatusForbidden,
			},
		},
		{
			name: "usecase unknown error",
			useCase: &useCaseMock{
				outDTO: nil,
				err:    errDefault,
			},
			want: want{
				err:  errDefault,
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

			h := NewHandler(tc.useCase).Handle

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
