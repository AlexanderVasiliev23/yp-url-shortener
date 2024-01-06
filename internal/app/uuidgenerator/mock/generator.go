package mock

import "github.com/google/uuid"

type Generator struct {
	res uuid.UUID
}

func NewGenerator(res uuid.UUID) *Generator {
	return &Generator{res: res}
}

func (g *Generator) Generate() uuid.UUID {
	return g.res
}
