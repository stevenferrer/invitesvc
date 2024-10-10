package postgres_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stevenferrer/invitesvc/postgres"
	"github.com/stevenferrer/invitesvc/postgres/txdb"
	"github.com/stevenferrer/invitesvc/token"
)

func TestTokenRepository(t *testing.T) {
	db := txdb.MustOpen()
	defer db.Close()

	// migrate db
	postgres.MustMigrate(db)

	ctx := context.TODO()
	tokenRepo := postgres.NewTokenRepository(db)

	t.Run("create and retrieve token", func(t *testing.T) {
		tokenID, err := token.NewID()
		require.NoError(t, err)

		err = tokenRepo.CreateToken(ctx, tokenID)
		require.NoError(t, err)

		gotToken, err := tokenRepo.GetToken(ctx, tokenID)
		require.NoError(t, err)

		assert.Equal(t, tokenID, gotToken.ID)
		assert.NotNil(t, gotToken.CreatedAt)
		assert.Nil(t, gotToken.RedeemedAt)

		t.Run("token not found", func(t *testing.T) {
			tokenID, err = token.NewID()
			require.NoError(t, err)

			gotToken, err = tokenRepo.GetToken(ctx, tokenID)
			assert.ErrorIs(t, err, token.ErrTokenNotFound)
			assert.Nil(t, gotToken)
		})
	})

	t.Run("list tokens", func(t *testing.T) {
		gotTokens, err := tokenRepo.ListTokens(ctx)
		require.NoError(t, err)
		assert.Len(t, gotTokens, 1)

	})

	t.Run("set token to disabled", func(t *testing.T) {
		gotTokens, err := tokenRepo.ListTokens(ctx)
		require.NoError(t, err)
		assert.Len(t, gotTokens, 1)

		tk := gotTokens[0]

		// set to disabled
		err = tokenRepo.SetTokenDisabled(ctx, tk.ID)
		require.NoError(t, err)

		// verify update
		gotTk, err := tokenRepo.GetToken(ctx, tk.ID)
		require.NoError(t, err)
		assert.True(t, gotTk.Disabled)

		t.Run("token not found", func(t *testing.T) {
			tokenID, err := token.NewID()
			require.NoError(t, err)

			err = tokenRepo.SetTokenDisabled(ctx, tokenID)
			assert.ErrorIs(t, err, token.ErrTokenNotFound)
		})
	})

	t.Run("set token to redeemed", func(t *testing.T) {
		gotTokens, err := tokenRepo.ListTokens(ctx)
		require.NoError(t, err)
		assert.Len(t, gotTokens, 1)

		tk := gotTokens[0]

		// set to redeemed
		err = tokenRepo.SetTokenRedeemed(ctx, tk.ID)
		require.NoError(t, err)

		// verify update
		gotTk, err := tokenRepo.GetToken(ctx, tk.ID)
		require.NoError(t, err)
		assert.NotNil(t, gotTk.RedeemedAt)

		t.Run("token not found", func(t *testing.T) {
			tokenID, err := token.NewID()
			require.NoError(t, err)

			err = tokenRepo.SetTokenDisabled(ctx, tokenID)
			assert.ErrorIs(t, err, token.ErrTokenNotFound)
		})
	})
}
