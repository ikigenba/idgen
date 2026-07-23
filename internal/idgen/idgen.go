// Package idgen provides pure identifier minting and decoding operations.
package idgen

import (
	"errors"
	"time"
)

// Clock supplies the current time to callers without coupling the core to the
// wall clock.
type Clock interface {
	Now() time.Time
}

// Epoch is the reference time used by the identifier encoding.
var Epoch time.Time

var errNotImplemented = errors.New("idgen: not implemented")

// MintAt mints n identifiers at t.
func MintAt(t time.Time, n int) (string, error) {
	_ = t
	_ = n
	return "", errNotImplemented
}

// TimeOf decodes the timestamp represented by id.
func TimeOf(id string) (time.Time, error) {
	_ = id
	return time.Time{}, errNotImplemented
}
