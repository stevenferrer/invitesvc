package authn

import (
	"context"

	"github.com/pkg/errors"
)

// Service is an auth service
type Service interface {
	// GenerateAuthKey generates an auth key
	GenerateAuthKey(context.Context) (AuthKey, error)
	// AuthKeyExists checks if auth key exists in db
	AuthKeyExists(context.Context, AuthKey) (bool, error)
}

// authService implements auth service
type authService struct {
	authRepo Repository
}

var _ Service = (*authService)(nil)

// NewAuthService takes an auth repo and returns an auth service
func NewAuthService(authRepo Repository) Service {
	return &authService{authRepo: authRepo}
}

func (svc *authService) GenerateAuthKey(ctx context.Context) (AuthKey, error) {
	authKey := NewAuthKey()
	err := svc.authRepo.CreateAuthKey(ctx, authKey)
	return authKey, errors.Wrap(err, "create auth key")
}

func (svc *authService) AuthKeyExists(ctx context.Context, authKey AuthKey) (bool, error) {
	return svc.authRepo.AuthKeyExists(ctx, authKey)
}
