package deleter

import (
	"context"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

type repositoryMock struct {
	err           error
	deletedTokens []string
}

func (m *repositoryMock) DeleteByTokens(ctx context.Context, tokens []string) error {
	m.deletedTokens = append(m.deletedTokens, tokens...)

	return m.err
}

func TestDeleteWorker_SettingOptions(t *testing.T) {
	maxBatchSize := 23
	savingInterval := time.Second
	repoDeletionTimeout := time.Minute

	worker := NewDeleteWorker(&repositoryMock{}, Options{
		MaxBatchSize:        maxBatchSize,
		SavingInterval:      savingInterval,
		RepoDeletionTimeout: repoDeletionTimeout,
	})

	assert.Equal(t, worker.maxBatchSize, maxBatchSize)
	assert.Equal(t, worker.savingInterval, savingInterval)
	assert.Equal(t, worker.repoDeletionTimeout, repoDeletionTimeout)
}

func BenchmarkTestDeleteWorker(b *testing.B) {
	testCases := []struct {
		name         string
		tokensToSave []string
		opts         Options
	}{
		{
			name:         "full batch saved",
			tokensToSave: []string{"token1", "token2"},
			opts: Options{
				MaxBatchSize:   2,
				SavingInterval: 1 * time.Second,
			},
		},
		{
			name:         "saved on tick",
			tokensToSave: []string{"token1"},
			opts: Options{
				MaxBatchSize:   2,
				SavingInterval: 1 * time.Millisecond,
			},
		},
		{
			name:         "saved on channel closing",
			tokensToSave: []string{"token1"},
			opts: Options{
				MaxBatchSize:   2,
				SavingInterval: 1 * time.Second,
			},
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			repo := &repositoryMock{}
			worker := NewDeleteWorker(repo, tc.opts)

			ch := make(chan DeleteTask, 1)

			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()

				for _, token := range tc.tokensToSave {
					ch <- DeleteTask{[]string{token}}
					time.Sleep(10 * time.Millisecond)
				}
				close(ch)
			}()

			worker.Consume(ch)

			wg.Wait()

			assert.Equal(b, tc.tokensToSave, repo.deletedTokens)
		})
	}
}
