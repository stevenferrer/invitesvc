package token

import (
	"context"

	"github.com/pkg/errors"
)

// Service is an invite service
type Service interface {
	// GenerateToken generates new invite token
	GenerateToken(context.Context) (ID, error)
	// GetToken retrieves an invite token
	GetToken(context.Context, ID) (*Token, error)
	// ListTokens retrieves the list of invite tokens
	ListTokens(context.Context) ([]*Token, error)
	// DisableToken is used to disable an invite token
	DisableToken(context.Context, ID) error
	// RedeemToken is used to redeem an invite token
	RedeemToken(context.Context, ID) error
}

// tokenService implements token service
type tokenService struct {
	repo Repository
}

var _ Service = (*tokenService)(nil)

// NewService returns a new token service
func NewService(tokenRepo Repository) Service {
	return &tokenService{repo: tokenRepo}
}

// GenerateToken generates a new token
func (svc *tokenService) GenerateToken(ctx context.Context) (ID, error) {
	id, err := NewID()
	if err != nil {
		return NilID, errors.Wrap(err, "new id")
	}

	// save token
	err = svc.repo.CreateToken(ctx, id)
	if err != nil {
		return NilID, errors.Wrap(err, "create token")
	}

	return id, nil
}

// GetToken retrieves a token
func (svc *tokenService) GetToken(ctx context.Context, id ID) (*Token, error) {
	token, err := svc.repo.GetToken(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "get token")
	}

	return token, nil
}

// ListTokens retrives list of tokens
func (svc *tokenService) ListTokens(ctx context.Context) ([]*Token, error) {
	tokens, err := svc.repo.ListTokens(ctx)
	return tokens, errors.Wrap(err, "list tokens")
}

// DisableToken disables a token
func (svc *tokenService) DisableToken(ctx context.Context, id ID) error {
	return svc.repo.SetTokenDisabled(ctx, id)
}

// RedeemToken redeems a token
func (svc *tokenService) RedeemToken(ctx context.Context, id ID) error {
	tk, err := svc.repo.GetToken(ctx, id)
	if err != nil {
		return errors.Wrap(err, "get token")
	}

	// validate token
	err = tk.Validate()
	if err != nil {
		return err
	}

	// redeem token
	err = svc.repo.SetTokenRedeemed(ctx, tk.ID)
	return errors.Wrap(err, "set token redeemed")
}
