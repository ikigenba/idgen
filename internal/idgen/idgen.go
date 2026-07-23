// Package idgen provides pure identifier minting and decoding operations.
package idgen

import (
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"
)

const (
	multiplier = int64(0x9E3779B1)
	offset     = int64(0xC0FFEE)
)

var (
	// Epoch is the zero point of the identifier timestamp.
	Epoch = time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)

	// ErrInvalidID is wrapped by errors returned for malformed identifiers.
	ErrInvalidID = errors.New("invalid id")

	modulus           = new(big.Int).Exp(big.NewInt(36), big.NewInt(8), nil)
	multiplierInverse *big.Int
	canonicalID       = regexp.MustCompile(`^[A-Za-z0-9]+-([0-9A-Z]{4})-([0-9A-Z]{4})$`)
)

func init() {
	multiplierInverse = mustModInverse(big.NewInt(multiplier), modulus)
}

func mustModInverse(value, modulus *big.Int) *big.Int {
	inverse := new(big.Int).ModInverse(value, modulus)
	if inverse == nil {
		panic("idgen: multiplier is not invertible modulo 36^8")
	}
	return inverse
}

// MintAt returns "<prefix>-XXXX-XXXX" for t. Times before Epoch are clamped
// to Epoch. Prefix is trusted to be a non-empty run of letters and digits.
func MintAt(prefix string, t time.Time) string {
	if t.Before(Epoch) {
		t = Epoch
	}

	ms := new(big.Int).Sub(
		big.NewInt(t.UnixMilli()),
		big.NewInt(Epoch.UnixMilli()),
	)
	encoded := new(big.Int).Mul(ms, big.NewInt(multiplier))
	encoded.Add(encoded, big.NewInt(offset))
	encoded.Mod(encoded, modulus)

	body := strings.ToUpper(encoded.Text(36))
	body = strings.Repeat("0", 8-len(body)) + body
	return prefix + "-" + body[:4] + "-" + body[4:]
}

// TimeOf returns the instant represented by id. The prefix is accepted but
// ignored because the timestamp is encoded entirely in the body.
func TimeOf(id string) (time.Time, error) {
	parts := canonicalID.FindStringSubmatch(id)
	if parts == nil {
		return time.Time{}, fmt.Errorf("%w: non-canonical format", ErrInvalidID)
	}

	body := parts[1] + parts[2]
	encoded, ok := new(big.Int).SetString(body, 36)
	if !ok {
		return time.Time{}, fmt.Errorf("%w: invalid base-36 body", ErrInvalidID)
	}

	ms := new(big.Int).Sub(encoded, big.NewInt(offset))
	ms.Mul(ms, multiplierInverse)
	ms.Mod(ms, modulus)
	return Epoch.Add(time.Duration(ms.Int64()) * time.Millisecond), nil
}
