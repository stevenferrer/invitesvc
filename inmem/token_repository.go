package inmem

import (
	"context"
	"time"

	"github.com/hashicorp/go-memdb"
	"github.com/pkg/errors"

	"github.com/stevenferrer/invitesvc/token"
)

// TokenRepository is an in-memory implementation of token.Repository
type TokenRepository struct {
	db *memdb.MemDB
}

var _ token.Repository = (*TokenRepository)(nil)

// NewTokenRepository returns a new token repository
func NewTokenRepository(db *memdb.MemDB) *TokenRepository {
	return &TokenRepository{db: db}
}

// CreateToken creates a new token and saves it to database
func (repo *TokenRepository) CreateToken(ctx context.Context, id token.ID) error {
	// write transactin
	txn := repo.db.Txn(true)
	now := time.Now()
	// insert token
	err := txn.Insert(tokensTable, &token.Token{
		ID:         id,
		CreatedAt:  &now,
		RedeemedAt: nil,
		Disabled:   false,
	})
	if err != nil {
		return errors.Wrap(err, "insert token")
	}

	// commit txn
	txn.Commit()
	return nil
}

// GetToken retrieves a token from the database
func (repo *TokenRepository) GetToken(ctx context.Context, id token.ID) (*token.Token, error) {
	// read-ony transaction
	txn := repo.db.Txn(false)
	defer txn.Abort()

	// retrieve token
	v, err := txn.First(tokensTable, "id", id)
	if err != nil {
		return nil, errors.Wrap(err, "get token")
	}

	// token not found
	if v == nil {
		return nil, token.ErrTokenNotFound
	}

	t, ok := v.(*token.Token)
	if !ok {
		return nil, errors.Errorf("unexpected value type %T, expecting %T", v, &token.Token{})
	}

	return t, nil
}

// ListToken retrieves a list of tokens from the database
func (repo *TokenRepository) ListTokens(ctx context.Context) ([]*token.Token, error) {
	txn := repo.db.Txn(false)
	defer txn.Abort()

	it, err := txn.Get(tokensTable, "id")
	if err != nil {
		return nil, errors.Wrap(err, "get tokens iterator")
	}

	tokens := make([]*token.Token, 0, 10)
	for v := it.Next(); v != nil; v = it.Next() {
		t, ok := v.(*token.Token)
		if !ok {
			return nil, errors.Errorf("unexpected value type %T, expecting %T", v, &token.Token{})
		}

		tokens = append(tokens, t)
	}

	return tokens, nil
}

// SetTokenDisabled sets a token to disabled
func (repo *TokenRepository) SetTokenDisabled(ctx context.Context, id token.ID) error {
	gotTk, err := repo.GetToken(ctx, id)
	if err != nil {
		return errors.Wrap(err, "get token")
	}

	newTk := &token.Token{
		ID:         gotTk.ID,
		Disabled:   true,
		RedeemedAt: gotTk.RedeemedAt,
		CreatedAt:  gotTk.CreatedAt,
	}

	txn := repo.db.Txn(true)
	err = txn.Insert(tokensTable, newTk)
	if err != nil {
		return errors.Wrap(err, "update token")
	}

	txn.Commit()
	return nil
}

// SetTokenRedeemed sets a token to redeemed
func (repo *TokenRepository) SetTokenRedeemed(ctx context.Context, id token.ID) error {
	gotTk, err := repo.GetToken(ctx, id)
	if err != nil {
		return errors.Wrap(err, "get token")
	}

	redeemedAt := time.Now()
	newTk := &token.Token{
		ID:         gotTk.ID,
		Disabled:   true,
		RedeemedAt: &redeemedAt,
		CreatedAt:  gotTk.CreatedAt,
	}

	txn := repo.db.Txn(true)
	err = txn.Insert(tokensTable, newTk)
	if err != nil {
		return errors.Wrap(err, "update token")
	}

	txn.Commit()
	return nil
}
