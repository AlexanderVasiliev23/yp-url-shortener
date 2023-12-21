package models

import "github.com/google/uuid"

type ShortLink struct {
	ID       string `json:"id"`
	Token    string `json:"token"`
	Original string `json:"original"`
	UserID   int    `json:"user_id"`
}

func NewShortLink(userID int, uuid uuid.UUID, token, original string) *ShortLink {
	return &ShortLink{
		ID:       uuid.String(),
		Token:    token,
		Original: original,
		UserID:   userID,
	}
}

func (l ShortLink) IsValid() bool {
	return l.ID != "" && l.Token != "" && l.Original != "" && l.UserID != 0
}
