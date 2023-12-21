package google

import "github.com/google/uuid"

type UUIDGenerator struct {
}

func (g UUIDGenerator) Generate() uuid.UUID {
	return uuid.New()
}
