package batch

import (
	"context"
	"errors"
	"fmt"
	"testing"

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
		err    error
		outDTO *OutDTO
	}

	testCases := []struct {
		userContextFetcher userContextFetcher
		tokenGenerator     tokenGenerator
		batchSaver         batchSaverMock
		name               string
		inDTO              InDTO
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
			inDTO: InDTO{
				Items: []InDTOItem{
					{
						CorrelationID: correlationID1,
						OriginalURL:   "https://test_url.com",
					},
					{
						CorrelationID: correlationID2,
						OriginalURL:   "https://test_url_2.com",
					},
				},
			},
			want: want{
				err: nil,
				outDTO: &OutDTO{
					Items: []OutDTOItem{
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
			},
		},
		{
			name:               "user fetching error",
			userContextFetcher: &mock.UserContextFetcherMock{Err: errDefault},
			batchSaver:         batchSaverMock{},
			inDTO: InDTO{
				Items: []InDTOItem{
					{
						CorrelationID: correlationID1,
						OriginalURL:   "https://test_url.com",
					},
				},
			},
			want: want{
				outDTO: nil,
				err:    errDefault,
			},
		},
		{
			name:               "token generator error",
			userContextFetcher: &mock.UserContextFetcherMock{},
			tokenGenerator:     tokenGeneratorMock{err: errDefault},
			batchSaver:         batchSaverMock{},
			inDTO: InDTO{
				Items: []InDTOItem{
					{
						CorrelationID: correlationID1,
						OriginalURL:   "https://test_url.com",
					},
				},
			},
			want: want{
				outDTO: nil,
				err:    errDefault,
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
			inDTO: InDTO{
				Items: []InDTOItem{
					{
						CorrelationID: correlationID1,
						OriginalURL:   "https://test_url.com",
					},
				},
			},
			want: want{
				outDTO: nil,
				err:    errDefault,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			useCase := NewUseCase(tc.batchSaver, tc.tokenGenerator, google.UUIDGenerator{}, tc.userContextFetcher, addr)

			outDTO, err := useCase.Shorten(context.Background(), tc.inDTO)

			if tc.want.err == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tc.want.err, err)
			}

			if tc.want.outDTO != nil {
				assert.Equal(t, *tc.want.outDTO, *outDTO)
			} else {
				assert.Nil(t, outDTO)
			}
		})
	}
}
