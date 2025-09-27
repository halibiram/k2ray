package utils

import (
	"crypto/rand"
	"math/big"

	"github.com/rs/zerolog/log"
)

// SecureIntn returns a cryptographically secure random integer in [0, n).
// It panics on error to simplify use in contexts where entropy is expected to be available.
func SecureIntn(n int64) int64 {
	if n <= 0 {
		// As this is for mock data, returning 0 is a safe default.
		return 0
	}
	result, err := rand.Int(rand.Reader, big.NewInt(n))
	if err != nil {
		// In a production system, you might want to handle this more gracefully,
		// but for a command-line tool or server startup, failing fast is okay.
		log.Fatal().Err(err).Msg("Failed to generate secure random number")
	}
	return result.Int64()
}

// SecureUint64n returns a cryptographically secure random unsigned integer in [0, n).
// It panics on error, similar to SecureIntn.
func SecureUint64n(n uint64) uint64 {
	if n == 0 {
		return 0
	}
	max := new(big.Int).SetUint64(n)
	result, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to generate secure random number")
	}
	return result.Uint64()
}

// SecureFloat64 returns a cryptographically secure random float64 in the range [0.0, 1.0).
func SecureFloat64() float64 {
	// The maximum value for the mantissa of a float64 is 2^53.
	// We generate a random integer up to this value and then divide to get a float.
	const maxMantissa = 1 << 53
	max := big.NewInt(maxMantissa)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to generate secure random number for float")
	}
	return float64(n.Int64()) / float64(maxMantissa)
}