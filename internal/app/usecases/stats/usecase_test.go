package stats

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	defaultUrlsCount           = 45
	defaultUsersCount          = 23
	defaultTrustedSubnet       = "127.0.0.1/24"
	defaultIPAddress           = "127.0.0.1"
	defaultNotTrustedIPAddress = "123.123.123.123"
)

var (
	errRepo = errors.New("repo error")
)

type repositoryMock struct {
	statusOutDTO *storage.StatsOutDTO
	err          error
}

func (r *repositoryMock) Stats(ctx context.Context) (*storage.StatsOutDTO, error) {
	return r.statusOutDTO, r.err
}

func TestHandle(t *testing.T) {
	type want struct {
		err    error
		outDTO *OutDTO
	}

	testCases := []struct {
		name string
		ip   string
		repo *repositoryMock
		want want
	}{
		{
			name: "success",
			ip:   defaultIPAddress,
			repo: &repositoryMock{
				statusOutDTO: &storage.StatsOutDTO{
					UrlsCount:  defaultUrlsCount,
					UsersCount: defaultUsersCount,
				},
				err: nil,
			},
			want: want{
				err: nil,
				outDTO: &OutDTO{
					Urls:  defaultUrlsCount,
					Users: defaultUsersCount,
				},
			},
		},
		{
			name: "invalid ip address",
			ip:   "invalid ip address",
			want: want{
				err:    fmt.Errorf("invalid ip address: %s", "invalid ip address"),
				outDTO: nil,
			},
		},
		{
			name: "not trusted ip address",
			ip:   defaultNotTrustedIPAddress,
			want: want{
				err:    ErrNotTrustedIP,
				outDTO: nil,
			},
		},
		{
			name: "repo error",
			ip:   defaultIPAddress,
			repo: &repositoryMock{
				err: errRepo,
			},
			want: want{
				err:    errRepo,
				outDTO: nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			useCase := NewUseCase(tc.repo, defaultTrustedSubnet)

			outDTO, err := useCase.Stats(context.Background(), tc.ip)

			if tc.want.err == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tc.want.err, err)
			}

			if tc.want.outDTO != nil {
				assert.Equal(t, tc.want.outDTO, outDTO)
			} else {
				assert.Nil(t, outDTO)
			}
		})
	}
}
