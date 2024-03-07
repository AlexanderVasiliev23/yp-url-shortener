// Package tokengenerator Пакет для генерации уникальных токенов
package tokengenerator

import (
	"errors"
	"math"
	"math/rand"
	"sync"
)

var symbols = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// ErrUniqueTokensRunOut missing godoc.
var ErrUniqueTokensRunOut = errors.New("unique tokens run out")

type token string

// TokenGenerator объект генератора токенов
type TokenGenerator struct {
	generatedTokens map[token]struct{}
	tokenLen        int
	maxTokensCount  int
	mu              sync.RWMutex
}

// New конструктор
func New(tokenLen int) *TokenGenerator {
	return &TokenGenerator{
		tokenLen:        tokenLen,
		maxTokensCount:  int(math.Pow(float64(len(symbols)), float64(tokenLen))),
		generatedTokens: make(map[token]struct{}),
	}
}

// Generate генерирует уникальный токен
func (g *TokenGenerator) Generate() (string, error) {
	if len(g.generatedTokens) >= g.maxTokensCount {
		return "", ErrUniqueTokensRunOut
	}

	for {
		t := g.generateRandom()

		if g.has(t) {
			continue
		}

		g.add(t)

		return string(t), nil
	}
}

func (g *TokenGenerator) has(t token) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	_, ok := g.generatedTokens[t]

	return ok
}

func (g *TokenGenerator) add(t token) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.generatedTokens[t] = struct{}{}
}

func (g *TokenGenerator) generateRandom() token {
	res := make([]rune, 0, g.tokenLen)

	for i := 0; i < g.tokenLen; i++ {
		pos := rand.Intn(len(symbols))
		res = append(res, symbols[pos])
	}

	return token(res)
}
