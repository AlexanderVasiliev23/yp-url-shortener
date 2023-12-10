package models

import "github.com/google/uuid"

type ShortLink struct {
	ID     string `json:"uuid"`
	Token  string `json:"short_url"`
	Origin string `json:"original_url"`
}

func NewShortLink(token, origin string) *ShortLink {
	return &ShortLink{
		ID:     uuid.NewString(),
		Token:  token,
		Origin: origin,
	}
}

func (l ShortLink) IsValid() bool {
	return l.ID != "" && l.Token != "" && l.Origin != ""
}
