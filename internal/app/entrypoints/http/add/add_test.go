package add

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/add"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	addr         = "localhost:8080"
	defaultToken = "default_test_token"
)

var (
	errDefault = errors.New("test_error")
)

type mockUseCase struct {
	err      error
	shortURL string
}

func (m *mockUseCase) Add(ctx context.Context, originalURL string) (shortenURL string, err error) {
	return m.shortURL, m.err
}

func TestAdd(t *testing.T) {
	type want struct {
		body string
		code int
	}

	tests := []struct {
		name    string
		method  string
		body    string
		want    want
		useCase useCase
	}{
		{
			name:   "success",
			method: http.MethodPost,
			body:   "test_url",
			useCase: &mockUseCase{
				err:      nil,
				shortURL: fmt.Sprintf("%s/%s", addr, defaultToken),
			},
			want: want{
				code: http.StatusCreated,
				body: fmt.Sprintf("%s/%s", addr, defaultToken),
			},
		},
		{
			name:   "empty body",
			method: http.MethodPost,
			body:   "",
			useCase: &mockUseCase{
				err:      add.ErrOriginalURLIsEmpty,
				shortURL: "",
			},
			want: want{
				code: http.StatusBadRequest,
				body: "",
			},
		},
		{
			name:   "usecase error",
			method: http.MethodPost,
			body:   "test_url",
			useCase: &mockUseCase{
				err:      errDefault,
				shortURL: "",
			},
			want: want{
				code: http.StatusInternalServerError,
				body: "",
			},
		},
		{
			name:   "already exists",
			method: http.MethodPost,
			body:   "test_url",
			useCase: &mockUseCase{
				err:      add.ErrOriginURLAlreadyExists,
				shortURL: fmt.Sprintf("%s/%s", addr, defaultToken),
			},
			want: want{
				code: http.StatusConflict,
				body: fmt.Sprintf("%s/%s", addr, defaultToken),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(tt.useCase).Add

			r := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(r, w)

			err := handler(c)

			if tt.want.code == http.StatusCreated {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.body, w.Body.String())
		})
	}
}
