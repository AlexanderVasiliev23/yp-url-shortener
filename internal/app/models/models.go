package models

import "github.com/google/uuid"

const (
	anonUserId = 999999
)

type ShortLink struct {
	ID       string `json:"id"`
	Token    string `json:"token"`
	Original string `json:"original"`
	UserId   int    `json:"user_id"`
}

func NewShortLink(userId int, uuid uuid.UUID, token, original string) *ShortLink {
	return &ShortLink{
		ID:       uuid.String(),
		Token:    token,
		Original: original,
		UserId:   userId,
	}
}

func NewShortLinkWithoutUserId(uuid uuid.UUID, token, original string) *ShortLink {
	return &ShortLink{
		ID:       uuid.String(),
		Token:    token,
		Original: original,
		UserId:   anonUserId,
	}
}

func (l ShortLink) IsValid() bool {
	return l.ID != "" && l.Token != "" && l.Original != "" && l.UserId != 0
}
