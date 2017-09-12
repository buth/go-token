// Package token implements tokens of arbitrary length that can be encoded to
// and from binary, text, and SQL database values. Generated text is safe for
// use in URLs.
package token

import (
	"crypto/rand"
	"database/sql/driver"
	"encoding/base64"
	"errors"
)

var (
	errInvalidTokenData = errors.New("token: invalid token data")
	tokenEncoding       = base64.URLEncoding.WithPadding(base64.NoPadding)
)

// Token is a slice of bytes.
type Token []byte

// NewToken returns a cryptographically secure pseudorandom Token of the given
// length (in bytes).
func New(length int) (Token, error) {
	t := make([]byte, length)
	if _, err := rand.Read(t); err != nil {
		return nil, err
	}

	return t, nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (t Token) MarshalBinary() ([]byte, error) {
	data := make([]byte, len(t))
	copy(data, t)
	return t, nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (t *Token) UnmarshalBinary(data []byte) error {
	*t = make([]byte, len(data))
	copy(*t, data)
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface, returning a
// base64-encoded slice of bytes.
func (t Token) MarshalText() (text []byte, err error) {
	data := make([]byte, tokenEncoding.EncodedLen(len(t)))
	tokenEncoding.Encode(data, t)
	return data, nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface, decoding a
// base64-encoded slice of bytes.
func (t *Token) UnmarshalText(text []byte) error {
	*t = make([]byte, tokenEncoding.DecodedLen(len(text)))
	_, err := tokenEncoding.Decode(*t, text)
	return err
}

// String returns a base64-encoded string.
func (t Token) String() string {
	return tokenEncoding.EncodeToString(t)
}

// Scan implements the sql.Scanner interface.
func (t *Token) Scan(src interface{}) error {
	data, ok := src.([]byte)
	if !ok {
		return errInvalidTokenData
	}

	return t.UnmarshalBinary(data)
}

// Value implements the driver.Valuer interface.
func (t Token) Value() (driver.Value, error) {
	return t.MarshalBinary()
}
