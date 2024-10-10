package token

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pkg/errors"
)

const (
	// alphabet is the token id alphabet
	alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	// idLen is the len of token id
	idLen = 12
	// tokenExpirDur is the token expiration duration
	tokenExpirDur = time.Hour * 24 * 7
)

// ID is a invite token id
type ID string

// NilID is a nil toke id
var NilID = ID("")

// NewID returns a new invite token id
func NewID() (ID, error) {
	// Refer to https://zelark.github.io/nano-id-cc/
	id, err := gonanoid.Generate(alphabet, idLen)
	if err != nil {
		return NilID, errors.Wrap(err, "generate id")
	}

	return ID(id), nil
}

// Token is an invite token
type Token struct {
	// ID is the token string
	ID ID `json:"id"`
	// Disabled is true when token is recalled/disabled
	Disabled bool `json:"disabled"`
	// RedeemedAt is the redeem timestamp
	RedeemedAt *time.Time `json:"redeemedAt"`
	// CreatedAt is the created at timestamp
	CreatedAt *time.Time `json:"createdAt"`
}

// Expiration returns the token expiration
func (t *Token) Expiration() time.Time {
	return t.CreatedAt.Add(tokenExpirDur)
}

// RedeemedAt returns true if the token is redeemed
func (t *Token) Redeemed() bool {
	return t.RedeemedAt != nil
}

// Expired returns true if the token is expired
func (t *Token) Expired() bool {
	return time.Now().After(t.Expiration())
}

// Validate token if possible to redeem
func (t *Token) Validate() error {
	// check if disabled
	if t.Disabled {
		return ErrTokenDisabled
	}

	// check if expired
	if t.Expired() {
		return ErrTokenExpired
	}

	// check if redeemed already
	if t.Redeemed() {
		return ErrTokenRedeemed
	}

	return nil
}
