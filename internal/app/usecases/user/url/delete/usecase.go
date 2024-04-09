package delete

import (
	"context"
	"errors"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/workers/deleter"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

type linksStorage interface {
	FilterOnlyThisUserTokens(ctx context.Context, userID int, tokens []string) ([]string, error)
}

type userContextFetcher interface {
	GetUserIDFromContext(ctx context.Context) (int, error)
}

type UseCase struct {
	linksStorage       linksStorage
	userContextFetcher userContextFetcher
	deleteByTokenCh    chan<- deleter.DeleteTask
}

func NewUseCase(linksStorage linksStorage, userContextFetcher userContextFetcher, deleteByTokenCh chan<- deleter.DeleteTask) *UseCase {
	return &UseCase{linksStorage: linksStorage, userContextFetcher: userContextFetcher, deleteByTokenCh: deleteByTokenCh}
}

func (u *UseCase) Delete(ctx context.Context, tokens []string) error {
	userID, err := u.userContextFetcher.GetUserIDFromContext(ctx)
	if err != nil {
		return ErrUnauthorized
	}

	thisUserTokens, err := u.linksStorage.FilterOnlyThisUserTokens(ctx, userID, tokens)
	if err != nil {
		return err
	}

	u.deleteByTokenCh <- deleter.DeleteTask{
		Tokens: thisUserTokens,
	}

	return nil
}
