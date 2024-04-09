package stats

import (
	"context"
	"errors"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
	iputil "github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/util/ip"
)

var (
	ErrNotTrustedIP = errors.New("not trusted IP")
)

type OutDTO struct {
	Urls  int
	Users int
}

type repository interface {
	Stats(ctx context.Context) (*storage.StatsOutDTO, error)
}

type UseCase struct {
	repository    repository
	trustedSubnet string
}

func NewUseCase(repository repository, trustedSubnet string) *UseCase {
	return &UseCase{repository: repository, trustedSubnet: trustedSubnet}
}

func (u *UseCase) Stats(ctx context.Context, ip string) (*OutDTO, error) {
	isTrusted, err := iputil.IsTrusted(ip, u.trustedSubnet)

	if err != nil {
		return nil, err
	}

	if !isTrusted {
		return nil, ErrNotTrustedIP
	}

	stats, err := u.repository.Stats(ctx)
	if err != nil {
		return nil, err
	}

	return &OutDTO{
		Urls:  stats.UrlsCount,
		Users: stats.UsersCount,
	}, nil
}
