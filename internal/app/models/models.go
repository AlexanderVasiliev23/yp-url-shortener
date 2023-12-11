package models

import "github.com/google/uuid"

type ShortLink struct {
	ID       string `json:"id"`
	Token    string `json:"token"`
	Original string `json:"original"`
}

func NewShortLink(token, original string) *ShortLink {
	return &ShortLink{
		ID:       uuid.NewString(),
		Token:    token,
		Original: original,
	}
}

func (l ShortLink) IsValid() bool {
	return l.ID != "" && l.Token != "" && l.Original != ""
}
