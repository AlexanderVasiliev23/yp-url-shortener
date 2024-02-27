package google

import "github.com/google/uuid"

// UUIDGenerator missing godoc.
type UUIDGenerator struct {
}

// Generate missing godoc.
func (g UUIDGenerator) Generate() uuid.UUID {
	return uuid.New()
}
