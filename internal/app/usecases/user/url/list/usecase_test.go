package list

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
)

const (
	defaultToken    = "test_token"
	defaultAddr     = "https://my_url_shortener"
	defaultUserID   = 1234
	defaultOriginal = "test_original"
)

var (
	ErrDefault = errors.New("test_error")
)

type userContextFetcherMock struct {
	err    error
	userID int
}

func (f userContextFetcherMock) GetUserIDFromContext(ctx context.Context) (int, error) {
	return f.userID, f.err
}

type storageMock struct {
	err    error
	result []*models.ShortLink
}

func (s storageMock) FindByUserID(ctx context.Context, userID int) ([]*models.ShortLink, error) {
	return s.result, s.err
}

func TestList(t *testing.T) {
	type want struct {
		err    error
		outDTO *OutDTO
	}

	testCases := []struct {
		userContextFetcherMock *userContextFetcherMock
		name                   string
		storage                storageMock
		want                   want
	}{
		{
			name: "success list",
			userContextFetcherMock: &userContextFetcherMock{
				userID: defaultUserID,
			},
			storage: storageMock{
				result: []*models.ShortLink{
					{
						Token:    defaultToken,
						Original: defaultOriginal,
					},
				},
			},
			want: want{
				err: nil,
				outDTO: &OutDTO{
					Items: []OutDTOItem{
						{
							ShortURL:    fmt.Sprintf("%s/%s", defaultAddr, defaultToken),
							OriginalURL: defaultOriginal,
						},
					},
				},
			},
		},
		{
			name: "empty list",
			userContextFetcherMock: &userContextFetcherMock{
				userID: defaultUserID,
			},
			storage: storageMock{
				result: make([]*models.ShortLink, 0),
			},
			want: want{
				err:    ErrNoSavedURLs,
				outDTO: nil,
			},
		},
		{
			name: "unauthorized",
			userContextFetcherMock: &userContextFetcherMock{
				err: ErrDefault,
			},
			want: want{
				err:    ErrUnauthorized,
				outDTO: nil,
			},
		},
		{
			name: "storage error",
			userContextFetcherMock: &userContextFetcherMock{
				userID: defaultUserID,
			},
			storage: storageMock{
				err: ErrDefault,
			},
			want: want{
				err:    ErrDefault,
				outDTO: nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			useCase := NewUseCase(tc.storage, tc.userContextFetcherMock, defaultAddr)

			outDTO, err := useCase.List(context.Background())

			if tc.want.err != nil {
				assert.Equal(t, tc.want.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if tc.want.outDTO != nil {
				assert.Equal(t, *tc.want.outDTO, *outDTO)
			} else {
				assert.Nil(t, outDTO)
			}
		})
	}
}
