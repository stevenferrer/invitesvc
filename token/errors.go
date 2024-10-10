package token

import (
	"github.com/pkg/errors"
)

// List of token related errors
var (
	ErrTokenNotFound = errors.New("token not found")
	ErrTokenDisabled = errors.New("token is disabled")
	ErrTokenExpired  = errors.New("token already expired")
	ErrTokenRedeemed = errors.New("token already redeemed")
)
