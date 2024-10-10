package token_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stevenferrer/invitesvc/postgres"
	"github.com/stevenferrer/invitesvc/postgres/txdb"
	"github.com/stevenferrer/invitesvc/token"
)

func TestService(t *testing.T) {
	db := txdb.MustOpen()
	defer db.Close()

	// migrate db
	postgres.MustMigrate(db)

	tokenRepo := postgres.NewTokenRepository(db)
	tokenSvc := token.NewService(tokenRepo)

	ctx := context.TODO()
	t.Run("generate and retrieve token", func(t *testing.T) {
		tokenID, err := tokenSvc.GenerateToken(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, tokenID)

		tk, err := tokenSvc.GetToken(ctx, tokenID)
		require.NoError(t, err)

		assert.Equal(t, tokenID, tk.ID)
		assert.False(t, tk.Disabled)
		assert.Nil(t, tk.RedeemedAt)
		assert.NotNil(t, tk.CreatedAt)

		t.Run("token not found", func(t *testing.T) {
			tokenID, err = token.NewID()
			require.NoError(t, err)

			_, err = tokenSvc.GetToken(ctx, tokenID)
			assert.ErrorIs(t, err, token.ErrTokenNotFound)
		})
	})

	t.Run("disable token", func(t *testing.T) {
		tokenID, err := tokenSvc.GenerateToken(ctx)
		require.NoError(t, err)

		// disable token
		err = tokenSvc.DisableToken(ctx, tokenID)
		require.NoError(t, err)

		// verify token is disabled
		gotTk, err := tokenSvc.GetToken(ctx, tokenID)
		require.NoError(t, err)
		assert.True(t, gotTk.Disabled)

		t.Run("token not found", func(t *testing.T) {
			tokenID, err = token.NewID()
			require.NoError(t, err)

			err = tokenSvc.DisableToken(ctx, tokenID)
			assert.ErrorIs(t, err, token.ErrTokenNotFound)
		})
	})

	t.Run("redeem token", func(t *testing.T) {
		tokenID, err := tokenSvc.GenerateToken(ctx)
		require.NoError(t, err)

		// redeem token
		err = tokenSvc.RedeemToken(ctx, tokenID)
		require.NoError(t, err)

		// verify token is redeemed
		gotTk, err := tokenSvc.GetToken(ctx, tokenID)
		require.NoError(t, err)
		assert.NotNil(t, gotTk.RedeemedAt)

		t.Run("token not found", func(t *testing.T) {
			tokenID, err = token.NewID()
			require.NoError(t, err)

			err = tokenSvc.RedeemToken(ctx, tokenID)
			assert.ErrorIs(t, err, token.ErrTokenNotFound)
		})
	})

	t.Run("list tokens", func(t *testing.T) {
		tokens, err := tokenSvc.ListTokens(ctx)
		require.NoError(t, err)
		assert.Len(t, tokens, 3)
	})
}
