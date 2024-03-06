package tokengenerator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenGenerator_Generate(t *testing.T) {
	const tokenLen = 8

	tokGen := New(tokenLen)

	token1, err := tokGen.Generate()
	require.NoError(t, err)
	token2, err := tokGen.Generate()
	require.NoError(t, err)

	assert.NotEqual(t, token1, token2)
}

func TestTokenGenerator_Generate_UniqueTokensRunOut(t *testing.T) {
	const tokenLen = 1

	tokGen := New(tokenLen)

	for i := 0; i < len(symbols); i++ {
		_, err := tokGen.Generate()
		require.NoError(t, err)
	}

	_, err := tokGen.Generate()
	assert.ErrorIs(t, ErrUniqueTokensRunOut, err)
}

func Benchmark_Generate(b *testing.B) {
	const tokenLen = 8

	tokGen := New(tokenLen)

	for i := 0; i < b.N; i++ {
		_, _ = tokGen.Generate()
	}
}
