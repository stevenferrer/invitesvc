package authn

import (
	"context"
)

// Repository is an auth repository
type Repository interface {
	// CreateAuthKey creates an auth key
	CreateAuthKey(context.Context, AuthKey) error
	// AuthKeyExists checks if auth key exists in db
	AuthKeyExists(context.Context, AuthKey) (bool, error)
}
