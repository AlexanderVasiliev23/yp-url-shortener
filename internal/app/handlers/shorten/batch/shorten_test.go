package batch

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/util/auth/mock"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/uuidgenerator/google"
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

type batchSaverMock struct {
	err error
}

func (m batchSaverMock) SaveBatch(ctx context.Context, shortLinks []*models.ShortLink) error {
	return m.err
}

type tokenGeneratorMock struct {
	tokensSeq chan string
	err       error
}

func (t tokenGeneratorMock) Generate() (string, error) {
	if t.err != nil {
		return "", t.err
	}
	return <-t.tokensSeq, nil
}

func TestShorten(t *testing.T) {
	type want struct {
		err  error
		body string
		code int
	}

	testCases := []struct {
		userContextFetcher userContextFetcher
		tokenGenerator     tokenGenerator
		batchSaver         batchSaverMock
		name               string
		reqBody            string
		want               want
	}{
		{
			name:               "success",
			userContextFetcher: &mock.UserContextFetcherMock{},
			tokenGenerator: func() tokenGenerator {
				tokens := []string{token1, token2}
				tokensChan := make(chan string, len(tokens))
				for _, token := range tokens {
					tokensChan <- token
				}
				return tokenGeneratorMock{tokensSeq: tokensChan}
			}(),
			batchSaver: batchSaverMock{},
			reqBody:    fmt.Sprintf(`[{"correlation_id": "%s","original_url": "https://test_url.com"},{"correlation_id": "%s","original_url": "https://test_url.com"}]`, correlationID1, correlationID2),
			want: want{
				code: http.StatusCreated,
				body: fmt.Sprintf(`[{"correlation_id":"%s","short_url":"%s/%s"},{"correlation_id":"%s","short_url":"%s/%s"}]`+"\n", correlationID1, addr, token1, correlationID2, addr, token2),
			},
		},
		{
			name:               "user fetching error",
			userContextFetcher: &mock.UserContextFetcherMock{Err: errDefault},
			batchSaver:         batchSaverMock{},
			reqBody:            fmt.Sprintf(`[{"correlation_id": "%s","original_url": "https://test_url.com"},{"correlation_id": "%s","original_url": "https://test_url.com"}]`, correlationID1, correlationID2),
			want: want{
				code: http.StatusInternalServerError,
				err:  errDefault,
			},
		},
		{
			name:               "token generator error",
			userContextFetcher: &mock.UserContextFetcherMock{},
			tokenGenerator:     tokenGeneratorMock{err: errDefault},
			batchSaver:         batchSaverMock{},
			reqBody:            fmt.Sprintf(`[{"correlation_id": "%s","original_url": "https://test_url.com"},{"correlation_id": "%s","original_url": "https://test_url.com"}]`, correlationID1, correlationID2),
			want: want{
				code: http.StatusInternalServerError,
				err:  errDefault,
			},
		},
		{
			name:               "batch saver error",
			userContextFetcher: &mock.UserContextFetcherMock{},
			tokenGenerator: func() tokenGenerator {
				tokens := []string{token1, token2}
				tokensChan := make(chan string, len(tokens))
				for _, token := range tokens {
					tokensChan <- token
				}
				return tokenGeneratorMock{tokensSeq: tokensChan}
			}(),
			batchSaver: batchSaverMock{err: errDefault},
			reqBody:    fmt.Sprintf(`[{"correlation_id": "%s","original_url": "https://test_url.com"},{"correlation_id": "%s","original_url": "https://test_url.com"}]`, correlationID1, correlationID2),
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

			h := NewShortener(tc.batchSaver, tc.tokenGenerator, google.UUIDGenerator{}, tc.userContextFetcher, addr).Handle

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
