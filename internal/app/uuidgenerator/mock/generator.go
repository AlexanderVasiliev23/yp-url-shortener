package mock

import "github.com/google/uuid"

// Generator missing godoc.
type Generator struct {
	res uuid.UUID
}

// NewGenerator missing godoc.
func NewGenerator(res uuid.UUID) *Generator {
	return &Generator{res: res}
}

// Generate missing godoc.
func (g *Generator) Generate() uuid.UUID {
	return g.res
}
