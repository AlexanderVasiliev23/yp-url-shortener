package deleteurl

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/util/auth/mock"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/workers/deleter"
)

var (
	errDefault = errors.New("default error")
)

type storageMock struct {
	result []string
	err    error
}

func (m storageMock) FilterOnlyThisUserTokens(ctx context.Context, userID int, tokens []string) ([]string, error) {
	return m.result, m.err
}

func TestDelete(t *testing.T) {
	type want struct {
		code  int
		err   error
		tasks []deleter.DeleteTask
	}

	testCases := []struct {
		name string
		body string
		want want

		userContextFetcher userContextFetcher
		linksStorage       storageMock
	}{
		{
			name: "success",
			body: `["token1", "token2"]`,
			want: want{
				code:  http.StatusAccepted,
				tasks: []deleter.DeleteTask{{Tokens: []string{"token1", "token2"}}},
			},
			userContextFetcher: &mock.UserContextFetcherMock{},
			linksStorage:       storageMock{result: []string{"token1", "token2"}},
		},
		{
			name: "bad request",
			body: "",
			want: want{
				code: http.StatusBadRequest,
				err:  io.EOF,
			},
			userContextFetcher: &mock.UserContextFetcherMock{},
			linksStorage:       storageMock{},
		},
		{
			name: "user fetcher error",
			body: `["token1", "token2"]`,
			want: want{
				code: http.StatusUnauthorized,
				err:  errDefault,
			},
			userContextFetcher: &mock.UserContextFetcherMock{Err: errDefault},
			linksStorage:       storageMock{},
		},
		{
			name: "repo error",
			body: `["token1", "token2"]`,
			want: want{
				code: http.StatusInternalServerError,
				err:  errDefault,
			},
			userContextFetcher: &mock.UserContextFetcherMock{},
			linksStorage:       storageMock{err: errDefault},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(tc.body))

			ch := make(chan deleter.DeleteTask, 1)

			h := NewHandler(tc.linksStorage, tc.userContextFetcher, ch).Delete

			e := echo.New()
			c := e.NewContext(req, recorder)

			err := h(c)
			close(ch)

			assert.ErrorIs(t, err, tc.want.err)
			assert.Equal(t, tc.want.code, recorder.Code)

			if len(tc.want.tasks) > 0 {
				chanAsSlice := chanToSlice(ch)
				assert.Equal(t, tc.want.tasks, chanAsSlice)
			}
		})
	}
}

func chanToSlice(ch chan deleter.DeleteTask) []deleter.DeleteTask {
	res := make([]deleter.DeleteTask, 0)

	for val := range ch {
		res = append(res, val)
	}

	return res
}
