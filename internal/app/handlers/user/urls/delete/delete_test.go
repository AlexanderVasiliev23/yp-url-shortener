package delete

import (
	"context"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type storageMock struct {
	result []*models.ShortLink
	err    error
}

func (m storageMock) DeleteTokens(ctx context.Context, userID int, tokens []string) error {
	return nil
}

type userContextFetcherMock struct {
	userID int
	err    error
}

func (f userContextFetcherMock) GetUserIDFromContext(ctx context.Context) (int, error) {
	return f.userID, f.err
}

func TestDelete(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(`["token1", "token2"]`))

	h := Delete(storageMock{}, userContextFetcherMock{})

	e := echo.New()
	c := e.NewContext(req, recorder)

	err := h(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, recorder.Code)
}
