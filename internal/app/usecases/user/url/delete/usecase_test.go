package delete

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/util/auth/mock"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/workers/deleter"
)

var (
	errDefault = errors.New("default error")
)

type storageMock struct {
	err    error
	result []string
}

func (m storageMock) FilterOnlyThisUserTokens(ctx context.Context, userID int, tokens []string) ([]string, error) {
	return m.result, m.err
}

func TestDelete(t *testing.T) {
	type want struct {
		err   error
		tasks []deleter.DeleteTask
	}

	testCases := []struct {
		linksStorage       storageMock
		userContextFetcher userContextFetcher
		name               string
		tokensToDelete     []string
		want               want
	}{
		{
			name:           "success",
			tokensToDelete: []string{"token1", "token2"},
			want: want{
				err:   nil,
				tasks: []deleter.DeleteTask{{Tokens: []string{"token1", "token2"}}},
			},
			userContextFetcher: &mock.UserContextFetcherMock{},
			linksStorage:       storageMock{result: []string{"token1", "token2"}},
		},
		{
			name:           "success: user owns only token2",
			tokensToDelete: []string{"token1", "token2"},
			want: want{
				err:   nil,
				tasks: []deleter.DeleteTask{{Tokens: []string{"token2"}}},
			},
			userContextFetcher: &mock.UserContextFetcherMock{},
			linksStorage:       storageMock{result: []string{"token2"}},
		},
		{
			name:           "user fetcher error",
			tokensToDelete: []string{"token1", "token2"},
			want: want{
				err: ErrUnauthorized,
			},
			userContextFetcher: &mock.UserContextFetcherMock{Err: errDefault},
			linksStorage:       storageMock{},
		},
		{
			name:           "repo error",
			tokensToDelete: []string{"token1", "token2"},
			want: want{
				err: errDefault,
			},
			userContextFetcher: &mock.UserContextFetcherMock{},
			linksStorage:       storageMock{err: errDefault},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ch := make(chan deleter.DeleteTask, 1)
			useCase := NewUseCase(tc.linksStorage, tc.userContextFetcher, ch)

			err := useCase.Delete(context.Background(), tc.tokensToDelete)

			if tc.want.err == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tc.want.err, err)
			}

			close(ch)

			if len(tc.want.tasks) > 0 {
				chanAsSlice := chanToSlice(ch)
				assert.Equal(t, tc.want.tasks, chanAsSlice)
			}
		})
	}
}

func chanToSlice(ch chan deleter.DeleteTask) []deleter.DeleteTask {
	res := make([]deleter.DeleteTask, 0)

	for val := range ch {
		res = append(res, val)
	}

	return res
}
