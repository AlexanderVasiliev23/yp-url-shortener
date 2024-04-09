package add

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/add"
	url_shortener "github.com/AlexanderVasiliev23/yp-url-shortener/proto/gen/proto"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
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
	tests := []struct {
		name        string
		method      string
		originalURL string
		want        url_shortener.AddResponse
		useCase     useCase
	}{
		{
			name:        "success",
			method:      http.MethodPost,
			originalURL: "test_url",
			useCase: &mockUseCase{
				err:      nil,
				shortURL: fmt.Sprintf("%s/%s", addr, defaultToken),
			},
			want: url_shortener.AddResponse{
				ShortURL: fmt.Sprintf("%s/%s", addr, defaultToken),
				Error:    "",
			},
		},
		{
			name:        "empty body",
			method:      http.MethodPost,
			originalURL: "",
			useCase: &mockUseCase{
				err:      add.ErrOriginalURLIsEmpty,
				shortURL: "",
			},
			want: url_shortener.AddResponse{
				ShortURL: "",
				Error:    add.ErrOriginalURLIsEmpty.Error(),
			},
		},
		{
			name:        "usecase error",
			method:      http.MethodPost,
			originalURL: "test_url",
			useCase: &mockUseCase{
				err:      errDefault,
				shortURL: "",
			},
			want: url_shortener.AddResponse{
				ShortURL: "",
				Error:    "unknown error",
			},
		},
		{
			name:        "already exists",
			method:      http.MethodPost,
			originalURL: "test_url",
			useCase: &mockUseCase{
				err:      add.ErrOriginURLAlreadyExists,
				shortURL: fmt.Sprintf("%s/%s", addr, defaultToken),
			},
			want: url_shortener.AddResponse{
				ShortURL: fmt.Sprintf("%s/%s", addr, defaultToken),
				Error:    add.ErrOriginURLAlreadyExists.Error(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewServer(tt.useCase)

			resp, err := server.Add(context.Background(), &url_shortener.AddRequest{
				OriginalURL: tt.originalURL,
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.want, *resp)
		})
	}
}
