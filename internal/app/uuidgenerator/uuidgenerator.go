package uuidgenerator

import "github.com/google/uuid"

type UUIDGenerator interface {
	Generate() uuid.UUID
}
