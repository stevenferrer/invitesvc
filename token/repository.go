package token

import "context"

// Repository is a token repository
type Repository interface {
	// CreateToken creates a token
	CreateToken(context.Context, ID) error
	// GetToken retrieves a token from db
	GetToken(context.Context, ID) (*Token, error)
	// ListTokens retrieves list of tokens from db
	ListTokens(context.Context) ([]*Token, error)
	// SetTokenDisabled sets a token to disabled
	SetTokenDisabled(context.Context, ID) error
	// SetTokenRedeemed sets a token to redeemed
	SetTokenRedeemed(context.Context, ID) error
}
