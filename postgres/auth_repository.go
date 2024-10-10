package postgres

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/stevenferrer/invitesvc/authn"
)

// AuthRepository is an auth repository that uses postgres as backend
type AuthRepository struct {
	db *sql.DB
}

var _ authn.Repository = (*AuthRepository)(nil)

// NewAuthRepository retuns a new auth repository
func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// CreateAuthKey inserts an auth key into the db
func (repo *AuthRepository) CreateAuthKey(ctx context.Context, authKey authn.AuthKey) error {
	stmnt := `insert into auth_keys (auth_key) values ($1)`
	_, err := repo.db.ExecContext(ctx, stmnt, authKey)
	return errors.Wrap(err, "insert auth key")
}

// AuthKeyExists checks db if auth key exists
func (repo *AuthRepository) AuthKeyExists(ctx context.Context, authKey authn.AuthKey) (bool, error) {
	stmnt := `select exists(select 1 from auth_keys where auth_key=$1)`
	var exists bool
	err := repo.db.QueryRowContext(ctx, stmnt, authKey).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, errors.Wrap(err, "query auth key exists")
	}

	return exists, nil
}
