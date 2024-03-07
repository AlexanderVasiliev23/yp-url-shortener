package uuidgenerator

import "github.com/google/uuid"

// UUIDGenerator missing godoc.
type UUIDGenerator interface {
	Generate() uuid.UUID
}
