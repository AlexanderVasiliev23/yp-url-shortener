package deleter

import (
	"context"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/logger"
	"time"
)

const (
	maxBatchSize   = 10
	savingInterval = 1 * time.Second
)

type Repository interface {
	DeleteByTokens(ctx context.Context, tokens []string) error
}

type DeleteTask struct {
	Tokens []string
}

type DeleteWorker struct {
	repo Repository
}

func NewDeleteWorker(repo Repository) *DeleteWorker {
	return &DeleteWorker{
		repo: repo,
	}
}

func (w DeleteWorker) Consume(ch <-chan DeleteTask) {
	ticker := time.NewTicker(savingInterval)

	batch := make([]string, 0, maxBatchSize)

	for {
		select {
		case <-ticker.C:
			if len(batch) == 0 {
				continue
			}
			if err := w.delete(batch); err != nil {
				logger.Log.Errorf("exec delete urls by tokens: %v", err)
			}
			batch = batch[:0]
		case task := <-ch:
			for _, token := range task.Tokens {
				batch = append(batch, token)
				if len(batch) < maxBatchSize {
					continue
				}
				if err := w.delete(batch); err != nil {
					logger.Log.Errorf("exec delete urls by tokens: %v", err)
				}
				batch = batch[:0]
			}
		default:
		}
	}
}

func (w DeleteWorker) delete(tokens []string) error {
	ctx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFn()

	if err := w.repo.DeleteByTokens(ctx, tokens); err != nil {
		return fmt.Errorf("exec delete by tokens query: %w", err)
	}

	return nil
}
