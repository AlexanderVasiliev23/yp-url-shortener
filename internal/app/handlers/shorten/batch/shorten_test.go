package batch

import (
	"context"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	correlationId1 = "correlationId1"
	correlationId2 = "correlationId2"
	token1         = "token1"
	token2         = "token2"
	addr           = "test_addr"
)

type batchSaverMock struct {
	err error
}

func (m batchSaverMock) SaveBatch(ctx context.Context, shortLinks []*models.ShortLink) error {
	return m.err
}

type tokenGeneratorMock struct {
	tokensSeq chan string
}

func (t tokenGeneratorMock) Generate() (string, error) {
	return <-t.tokensSeq, nil
}

func TestShorten(t *testing.T) {
	reqBody := fmt.Sprintf(`[{"correlation_id": "%s","original_url": "https://test_url.com"},{"correlation_id": "%s","original_url": "https://test_url.com"}]`, correlationId1, correlationId2)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(reqBody))
	resp := httptest.NewRecorder()

	mockBatchSaver := batchSaverMock{}
	tokens := []string{token1, token2}
	tokensChan := make(chan string, len(tokens))
	for _, token := range tokens {
		tokensChan <- token
	}
	tokenGeneratorMock := tokenGeneratorMock{tokensSeq: tokensChan}
	h := Shorten(mockBatchSaver, tokenGeneratorMock, addr)

	e := echo.New()
	c := e.NewContext(req, resp)

	err := h(c)

	assert.NoError(t, err)
	expectedBody := fmt.Sprintf(`[{"correlation_id":"%s","short_url":"%s/%s"},{"correlation_id":"%s","short_url":"%s/%s"}]`+"\n", correlationId1, addr, token1, correlationId2, addr, token2)
	assert.Equal(t, expectedBody, resp.Body.String())
}
