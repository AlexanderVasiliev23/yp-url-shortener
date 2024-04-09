package add

import (
	"context"
	"errors"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/add"
	url_shortener "github.com/AlexanderVasiliev23/yp-url-shortener/proto/gen/proto"
)

type useCase interface {
	Add(ctx context.Context, originalURL string) (shortURL string, err error)
}

type Server struct {
	url_shortener.UnimplementedUrlShortenerServer

	useCase useCase
}

func NewServer(useCase useCase) *Server {
	return &Server{useCase: useCase}
}

func (s *Server) Add(ctx context.Context, request *url_shortener.AddRequest) (*url_shortener.AddResponse, error) {
	shortURL, err := s.useCase.Add(ctx, request.OriginalURL)

	if err != nil {
		if errors.Is(err, add.ErrOriginURLAlreadyExists) {
			return &url_shortener.AddResponse{
				ShortURL: shortURL,
				Error:    add.ErrOriginURLAlreadyExists.Error(),
			}, nil
		}

		if errors.Is(err, add.ErrOriginalURLIsEmpty) {
			return &url_shortener.AddResponse{
				ShortURL: "",
				Error:    add.ErrOriginalURLIsEmpty.Error(),
			}, nil
		}

		return &url_shortener.AddResponse{
			ShortURL: "",
			Error:    "unknown error",
		}, nil
	}

	return &url_shortener.AddResponse{
		ShortURL: shortURL,
	}, nil
}
