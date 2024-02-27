package models

import (
	"time"

	"github.com/google/uuid"
)

// ShortLink missing godoc.
type ShortLink struct {
	ID        string `json:"id"`
	Token     string `json:"token"`
	Original  string `json:"original"`
	UserID    int    `json:"user_id"`
	DeletedAt *time.Time
}

// NewShortLink missing godoc.
func NewShortLink(userID int, uuid uuid.UUID, token, original string) *ShortLink {
	return &ShortLink{
		ID:       uuid.String(),
		Token:    token,
		Original: original,
		UserID:   userID,
	}
}

// IsValid missing godoc.
func (l *ShortLink) IsValid() bool {
	return l.ID != "" && l.Token != "" && l.Original != "" && l.UserID != 0
}

// Delete missing godoc.
func (l *ShortLink) Delete() {
	at := time.Now()
	l.DeletedAt = &at
}
