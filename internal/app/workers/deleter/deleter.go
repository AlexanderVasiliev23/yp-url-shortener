package deleter

import (
	"context"
	"time"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/logger"
)

const (
	defaultMaxBatchSize        = 10
	defaultSavingInterval      = 1 * time.Second
	defaultRepoDeletionTimeout = 30 * time.Second
)

// Repository missing godoc.
type Repository interface {
	DeleteByTokens(ctx context.Context, tokens []string) error
}

// DeleteTask missing godoc.
type DeleteTask struct {
	Tokens []string
}

// DeleteWorker missing godoc.
type DeleteWorker struct {
	repo                Repository
	maxBatchSize        int
	savingInterval      time.Duration
	repoDeletionTimeout time.Duration
}

// Options missing godoc.
type Options struct {
	MaxBatchSize        int
	SavingInterval      time.Duration
	RepoDeletionTimeout time.Duration
}

// NewDeleteWorker missing godoc.
func NewDeleteWorker(repo Repository, opts Options) *DeleteWorker {
	worker := &DeleteWorker{
		repo:                repo,
		maxBatchSize:        defaultMaxBatchSize,
		savingInterval:      defaultSavingInterval,
		repoDeletionTimeout: defaultRepoDeletionTimeout,
	}

	if opts.MaxBatchSize != 0 {
		worker.maxBatchSize = opts.MaxBatchSize
	}

	if opts.SavingInterval != 0 {
		worker.savingInterval = opts.SavingInterval
	}

	if opts.RepoDeletionTimeout != 0 {
		worker.repoDeletionTimeout = opts.RepoDeletionTimeout
	}

	return worker
}

// Consume missing godoc.
func (w DeleteWorker) Consume(ch <-chan DeleteTask) {
	ticker := time.NewTicker(w.savingInterval)

	batch := make([]string, 0, w.maxBatchSize)

	for {
		select {
		case <-ticker.C:
			if len(batch) == 0 {
				continue
			}
			w.delete(batch)
			batch = batch[:0]
		case task, ok := <-ch:
			if !ok {
				if len(batch) > 0 {
					w.delete(batch)
				}
				return
			}
			for _, token := range task.Tokens {
				batch = append(batch, token)
				if len(batch) < w.maxBatchSize {
					continue
				}
				w.delete(batch)
				batch = batch[:0]
			}
		}
	}
}

func (w DeleteWorker) delete(tokens []string) {
	ctx, cancelFn := context.WithTimeout(context.Background(), w.repoDeletionTimeout)
	defer cancelFn()

	if err := w.repo.DeleteByTokens(ctx, tokens); err != nil {
		logger.Log.Errorf("exec delete urls by tokens: %v", err)
	}
}
