package batch

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/shorten/batch"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const (
	correlationID1 = "correlationId1"
	correlationID2 = "correlationId2"
	token1         = "token1"
	token2         = "token2"
	addr           = "https://my_url_shortener"
)

var (
	errDefault = errors.New("default_err")
)

type useCaseMock struct {
	out *batch.OutDTO
	err error
}

func (m *useCaseMock) Shorten(ctx context.Context, in batch.InDTO) (*batch.OutDTO, error) {
	return m.out, m.err
}

func TestShorten(t *testing.T) {
	type want struct {
		err  error
		body string
		code int
	}

	testCases := []struct {
		useCase useCase
		name    string
		reqBody string
		want    want
	}{
		{
			name: "success",
			useCase: &useCaseMock{
				out: &batch.OutDTO{
					Items: []batch.OutDTOItem{
						{
							CorrelationID: correlationID1,
							ShortURL:      fmt.Sprintf("%s/%s", addr, token1),
						},
						{
							CorrelationID: correlationID2,
							ShortURL:      fmt.Sprintf("%s/%s", addr, token2),
						},
					},
				},
				err: nil,
			},
			reqBody: fmt.Sprintf(`[{"correlation_id": "%s","original_url": "https://test_url.com"},{"correlation_id": "%s","original_url": "https://test_url.com"}]`, correlationID1, correlationID2),
			want: want{
				code: http.StatusCreated,
				body: fmt.Sprintf(`[{"correlation_id":"%s","short_url":"%s/%s"},{"correlation_id":"%s","short_url":"%s/%s"}]`+"\n", correlationID1, addr, token1, correlationID2, addr, token2),
			},
		},
		{
			name: "user fetching error",
			useCase: &useCaseMock{
				err: errDefault,
				out: nil,
			},
			reqBody: fmt.Sprintf(`[{"correlation_id": "%s","original_url": "https://test_url.com"},{"correlation_id": "%s","original_url": "https://test_url.com"}]`, correlationID1, correlationID2),
			want: want{
				code: http.StatusInternalServerError,
				err:  errDefault,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.reqBody))
			resp := httptest.NewRecorder()

			h := NewShortener(tc.useCase).Handle

			e := echo.New()
			c := e.NewContext(req, resp)

			err := h(c)

			if tc.want.err != nil {
				assert.ErrorIs(t, err, tc.want.err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.want.body, resp.Body.String())
		})
	}
}
