package deleteurl

import (
	"context"
	"errors"
	deleteusecase "github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/user/url/delete"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/workers/deleter"
)

var (
	errDefault = errors.New("default error")
)

type useCaseMock struct {
	err error
}

func (u *useCaseMock) Delete(ctx context.Context, tokens []string) error {
	return u.err
}

func TestDelete(t *testing.T) {
	type want struct {
		err  error
		code int
	}

	testCases := []struct {
		useCase useCase
		name    string
		body    string
		want    want
	}{
		{
			name: "success",
			body: `["token1", "token2"]`,
			want: want{
				code: http.StatusAccepted,
			},
			useCase: &useCaseMock{err: nil},
		},
		{
			name: "bad request",
			body: "",
			want: want{
				code: http.StatusBadRequest,
				err:  io.EOF,
			},
		},
		{
			name: "usecase unknown error",
			body: `["token1", "token2"]`,
			want: want{
				code: http.StatusInternalServerError,
				err:  errDefault,
			},
			useCase: &useCaseMock{err: errDefault},
		},
		{
			name: "unauthorized",
			body: `["token1", "token2"]`,
			want: want{
				code: http.StatusUnauthorized,
				err:  deleteusecase.ErrUnauthorized,
			},
			useCase: &useCaseMock{err: deleteusecase.ErrUnauthorized},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(tc.body))

			h := NewHandler(tc.useCase).Delete

			e := echo.New()
			c := e.NewContext(req, recorder)

			err := h(c)

			if tc.want.err != nil {
				assert.Equal(t, tc.want.err.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tc.want.code, recorder.Code)
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
