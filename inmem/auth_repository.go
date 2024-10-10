package inmem

import (
	"context"
	"time"

	"github.com/stevenferrer/invitesvc/authn"

	"github.com/hashicorp/go-memdb"
	"github.com/pkg/errors"
)

// AuthRepository is an in-memory implementation of authn.Repository
type AuthRepository struct {
	db *memdb.MemDB
}

var _ authn.Repository = (*AuthRepository)(nil)

// NewAuthRepository retuns a new auth repository
func NewAuthRepository(db *memdb.MemDB) *AuthRepository {
	return &AuthRepository{db: db}
}

// CreateAuthKey inserts an auth key into the db
func (repo *AuthRepository) CreateAuthKey(ctx context.Context, authKey authn.AuthKey) error {
	txn := repo.db.Txn(true)
	now := time.Now()
	err := txn.Insert(authsTable, &authn.Auth{
		Auth:      authKey,
		CreatedAt: &now,
	})
	if err != nil {
		return errors.Wrap(err, "insert auth key")
	}

	txn.Commit()
	return nil
}

// AuthKeyExists checks db if auth key exists
func (repo *AuthRepository) AuthKeyExists(ctx context.Context, authKey authn.AuthKey) (bool, error) {
	txn := repo.db.Txn(false)
	defer txn.Abort()

	v, err := txn.First(authsTable, "id", authKey)
	if err != nil {
		return false, err
	}

	// auth key not found
	if v == nil {
		return false, nil
	}

	_, ok := v.(*authn.Auth)
	if !ok {
		return false, errors.Errorf("unexpected value type %T, expecting %T", v, &authn.Auth{})
	}

	return true, nil
}
