package postgres

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/stevenferrer/invitesvc/token"
)

// TokenRepository as a token repository that uses postgres as backend
type TokenRepository struct {
	db *sql.DB
}

var _ token.Repository = (*TokenRepository)(nil)

// NewTokenRepository returns a token repository
func NewTokenRepository(db *sql.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

// CreateToken creates a new token and saves it to database
func (repo *TokenRepository) CreateToken(ctx context.Context, id token.ID) error {
	stmnt := `insert into tokens (token) values ($1)`
	_, err := repo.db.ExecContext(ctx, stmnt, id)
	return errors.Wrap(err, "insert token")
}

// GetToken retrieves a token from the database
func (repo *TokenRepository) GetToken(ctx context.Context, id token.ID) (*token.Token, error) {
	stmnt := `select token, disabled, redeemed_at, 
		created_at from tokens where token = $1`
	var tk token.Token
	err := repo.db.QueryRowContext(ctx, stmnt, id).Scan(
		&tk.ID, &tk.Disabled, &tk.RedeemedAt, &tk.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, token.ErrTokenNotFound
		}

		return nil, errors.Wrap(err, "query token")
	}

	return &tk, nil
}

// ListToken retrieves tokens from the database
func (repo *TokenRepository) ListTokens(ctx context.Context) ([]*token.Token, error) {
	stmnt := `select token, disabled, redeemed_at, created_at 
		from tokens order by created_at`
	rows, err := repo.db.QueryContext(ctx, stmnt)
	if err != nil {
		return nil, errors.Wrap(err, "query tokens")
	}

	tokens := make([]*token.Token, 0, 10)
	for rows.Next() {
		var tk token.Token
		err = rows.Scan(&tk.ID, &tk.Disabled, &tk.RedeemedAt, &tk.CreatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "scan row")
		}

		tokens = append(tokens, &tk)
	}

	return tokens, nil
}

// SetTokenDisabled sets a token to disabled
func (repo *TokenRepository) SetTokenDisabled(ctx context.Context, id token.ID) error {
	tk, err := repo.GetToken(ctx, id)
	if err != nil {
		return errors.Wrap(err, "get token")
	}

	stmnt := `update tokens set disabled=TRUE, 
		updated_at=now() where token=$1`
	_, err = repo.db.ExecContext(ctx, stmnt, tk.ID)
	return errors.Wrap(err, "update token")
}

// SetTokenRedeemed sets a token to redeemed
func (repo *TokenRepository) SetTokenRedeemed(ctx context.Context, id token.ID) error {
	tk, err := repo.GetToken(ctx, id)
	if err != nil {
		return errors.Wrap(err, "get token")
	}

	stmnt := `update tokens set redeemed_at=now(), 
		updated_at=now() where token=$1`
	_, err = repo.db.ExecContext(ctx, stmnt, tk.ID)
	return errors.Wrap(err, "update token")
}
