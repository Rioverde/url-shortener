package service

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
)

// alphabet is the set of characters used to build short keys.
const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// minKeyLength is the minimum allowed key length.
// Anything shorter would be too easy to brute-force.
const minKeyLength = 5

type CodeGenerator interface {
	GenerateRandomString(n int) (string, error)
}

type cryptoGenerator struct {
	// RandReader is the entropy source used by GenerateRandomString.
	// It defaults to crypto/rand.Reader and may be swapped in tests
	// to make output deterministic.
	rand io.Reader
}

func NewCryptoGenerator(r io.Reader) CodeGenerator {
	return &cryptoGenerator{
		rand: r,
	}
}

// GenerateRandomString returns a random string of length n drawn from alphabet.
// It returns an error if n is less than minKeyLength,
// or if the underlying entropy source fails.
func (g *cryptoGenerator) GenerateRandomString(n int) (string, error) {
	if n < minKeyLength {
		return "", fmt.Errorf("%w: must be at least %d, got %d", ErrKeyTooShort, minKeyLength, n)
	}
	// alphabetSize is the upper bound for picking an index into alphabet.
	// We hoist it out of the loop so big.NewInt is allocated only once.
	alphabetSize := big.NewInt(int64(len(alphabet)))
	// Pre-allocate the result slice with the exact capacity we need.
	m := make([]byte, n)
	// Pick a random character from the alphabet for each position in the slice.
	for i := range n {
		// Random index into the alphabet, drawn from a cryptographic source.
		val, err := rand.Int(g.rand, alphabetSize)
		if err != nil {
			return "", fmt.Errorf("get random value: %w", err)
		}
		// Place the chosen character into the result.
		m[i] = alphabet[val.Int64()]
	}
	return string(m), nil
}
