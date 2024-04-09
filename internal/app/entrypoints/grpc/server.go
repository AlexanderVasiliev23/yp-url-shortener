package grpc

import (
	"context"
	"errors"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/add"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/shorten/batch"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/shorten/single"
	url_shortener "github.com/AlexanderVasiliev23/yp-url-shortener/proto/gen/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type addUseCase interface {
	Add(ctx context.Context, originalURL string) (shortURL string, err error)
}

type batchUseCase interface {
	Shorten(ctx context.Context, in batch.InDTO) (*batch.OutDTO, error)
}

type singleUseCase interface {
	Shorten(ctx context.Context, originalURL string) (shortURL string, err error)
}

type Server struct {
	url_shortener.UnimplementedUrlShortenerServer

	addUseCase    addUseCase
	batchUseCase  batchUseCase
	singleUseCase singleUseCase
}

func NewServer(addUseCase addUseCase, batchUseCase batchUseCase, singleUseCase singleUseCase) *Server {
	return &Server{addUseCase: addUseCase, batchUseCase: batchUseCase, singleUseCase: singleUseCase}
}

func (s *Server) Add(ctx context.Context, request *url_shortener.AddRequest) (*url_shortener.AddResponse, error) {
	shortURL, err := s.addUseCase.Add(ctx, request.OriginalURL)

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

// CreateBatch missing godoc.
func (s *Server) CreateBatch(ctx context.Context, req *url_shortener.CreateBatchRequest) (*url_shortener.CreateBatchResponse, error) {
	inDTO := batch.InDTO{Items: make([]batch.InDTOItem, 0, len(req.GetItems()))}
	for _, item := range req.GetItems() {
		inDTO.Items = append(inDTO.Items, batch.InDTOItem{
			CorrelationID: item.CorrelationID,
			OriginalURL:   item.OriginalURL,
		})
	}

	outDTO, err := s.batchUseCase.Shorten(ctx, inDTO)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to shorten: %v", err)
	}

	resp := &url_shortener.CreateBatchResponse{
		Items: make([]*url_shortener.CreateBatchResponse_Item, 0, len(outDTO.Items)),
	}

	for _, item := range outDTO.Items {
		resp.Items = append(resp.Items, &url_shortener.CreateBatchResponse_Item{
			CorrelationID: item.CorrelationID,
			ShortURL:      item.ShortURL,
		})
	}

	return resp, nil
}

func (s *Server) CreateSingle(ctx context.Context, req *url_shortener.CreateSingleRequest) (*url_shortener.CreateSingleResponse, error) {
	shortURL, err := s.singleUseCase.Shorten(ctx, req.GetOriginalURL())

	if err != nil {
		if errors.Is(err, single.ErrEmptyOriginalURL) {
			return nil, status.Error(codes.InvalidArgument, "empty origin URL")
		}

		if errors.Is(err, single.ErrAlreadyExists) {
			return &url_shortener.CreateSingleResponse{
				ShortURL: shortURL,
			}, status.Error(codes.AlreadyExists, "already exists")
		}

		return nil, status.Error(codes.Internal, "unknown error")
	}

	return &url_shortener.CreateSingleResponse{
		ShortURL: shortURL,
	}, nil
}
