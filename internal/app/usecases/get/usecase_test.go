package get

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage/local"

	"github.com/stretchr/testify/assert"
)

const (
	defaultToken       = "default_test_token"
	defaultOriginalURL = "default_saved_url"
	defaultUserID      = 123
)

type mockRepo struct {
	err error
	url *models.ShortLink
}

func (m mockRepo) Get(ctx context.Context, s string) (*models.ShortLink, error) {
	return m.url, m.err
}

func TestGet(t *testing.T) {
	type want struct {
		err         error
		originalURL string
	}

	testCases := []struct {
		name  string
		repo  mockRepo
		token string
		want  want
	}{
		{
			name:  "success",
			repo:  mockRepo{url: models.NewShortLink(defaultUserID, uuid.New(), defaultToken, defaultOriginalURL)},
			token: defaultToken,
			want: want{
				err:         nil,
				originalURL: defaultOriginalURL,
			},
		},
		{
			name:  "token not found in repo",
			repo:  mockRepo{err: local.ErrURLNotFound},
			token: defaultToken,
			want: want{
				err:         local.ErrURLNotFound,
				originalURL: "",
			},
		},
		{
			name:  "empty token",
			repo:  mockRepo{},
			token: "",
			want: want{
				err:         ErrTokenIsEmpty,
				originalURL: "",
			},
		},
		{
			name: "deleted url",
			repo: mockRepo{url: func() *models.ShortLink {
				m := models.NewShortLink(defaultUserID, uuid.New(), defaultToken, defaultOriginalURL)
				m.Delete()
				return m
			}()},
			token: defaultToken,
			want: want{
				err:         ErrTokenIsDeleted,
				originalURL: "",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			useCase := NewUseCase(tc.repo)
			originalURL, err := useCase.Get(context.Background(), tc.token)

			if tc.want.err != nil {
				assert.Equal(t, tc.want.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.want.originalURL, originalURL)
		})
	}
}
