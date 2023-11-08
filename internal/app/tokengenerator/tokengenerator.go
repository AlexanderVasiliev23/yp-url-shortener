package tokengenerator

import (
	"errors"
	"math"
	"math/rand"
)

var symbols = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var ErrUniqueTokensRunOut = errors.New("unique tokens run out")

type (
	token  string
	tokens map[token]struct{}

	TokenGenerator struct {
		generatedTokens tokens
		tokenLen        int
		maxTokensCount  int
	}
)

func (m tokens) has(t token) bool {
	_, ok := m[t]

	return ok
}

func (m tokens) add(t token) {
	m[t] = struct{}{}
}

func New(tokenLen int) *TokenGenerator {
	return &TokenGenerator{
		tokenLen:        tokenLen,
		maxTokensCount:  int(math.Pow(float64(len(symbols)), float64(tokenLen))),
		generatedTokens: make(map[token]struct{}),
	}
}

func (g *TokenGenerator) Generate() (string, error) {
	if len(g.generatedTokens) >= g.maxTokensCount {
		return "", ErrUniqueTokensRunOut
	}

	for {
		t := g.generateRandom()

		if g.generatedTokens.has(t) {
			continue
		}

		g.generatedTokens.add(t)

		return string(t), nil
	}
}

func (g *TokenGenerator) generateRandom() token {
	res := make([]rune, 0, g.tokenLen)

	for i := 0; i < g.tokenLen; i++ {
		pos := rand.Intn(len(symbols))
		res = append(res, symbols[pos])
	}

	return token(res)
}
