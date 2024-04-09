package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/add"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/shorten/batch"
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

type mockAddUseCase struct {
	err      error
	shortURL string
}

func (m *mockAddUseCase) Add(ctx context.Context, originalURL string) (shortenURL string, err error) {
	return m.shortURL, m.err
}

type mockBatchUseCase struct {
	err    error
	outDTO *batch.OutDTO
}

func (m *mockBatchUseCase) Shorten(ctx context.Context, in batch.InDTO) (*batch.OutDTO, error) {
	return m.outDTO, m.err
}

type mockSingleUseCase struct {
	err      error
	shortURL string
}

func (m *mockSingleUseCase) Shorten(ctx context.Context, originalURL string) (shortURL string, err error) {
	return m.shortURL, m.err
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		originalURL string
		want        *url_shortener.AddResponse
		useCase     addUseCase
	}{
		{
			name:        "success",
			method:      http.MethodPost,
			originalURL: "test_url",
			useCase: &mockAddUseCase{
				err:      nil,
				shortURL: fmt.Sprintf("%s/%s", addr, defaultToken),
			},
			want: &url_shortener.AddResponse{
				ShortURL: fmt.Sprintf("%s/%s", addr, defaultToken),
				Error:    "",
			},
		},
		{
			name:        "empty body",
			method:      http.MethodPost,
			originalURL: "",
			useCase: &mockAddUseCase{
				err:      add.ErrOriginalURLIsEmpty,
				shortURL: "",
			},
			want: &url_shortener.AddResponse{
				ShortURL: "",
				Error:    add.ErrOriginalURLIsEmpty.Error(),
			},
		},
		{
			name:        "usecase error",
			method:      http.MethodPost,
			originalURL: "test_url",
			useCase: &mockAddUseCase{
				err:      errDefault,
				shortURL: "",
			},
			want: &url_shortener.AddResponse{
				ShortURL: "",
				Error:    "unknown error",
			},
		},
		{
			name:        "already exists",
			method:      http.MethodPost,
			originalURL: "test_url",
			useCase: &mockAddUseCase{
				err:      add.ErrOriginURLAlreadyExists,
				shortURL: fmt.Sprintf("%s/%s", addr, defaultToken),
			},
			want: &url_shortener.AddResponse{
				ShortURL: fmt.Sprintf("%s/%s", addr, defaultToken),
				Error:    add.ErrOriginURLAlreadyExists.Error(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewServer(tt.useCase, &mockBatchUseCase{}, &mockSingleUseCase{})

			resp, err := server.Add(context.Background(), &url_shortener.AddRequest{
				OriginalURL: tt.originalURL,
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.want, resp)
		})
	}
}
