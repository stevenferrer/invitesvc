package postgres_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stevenferrer/invitesvc/authn"
	"github.com/stevenferrer/invitesvc/postgres"
	"github.com/stevenferrer/invitesvc/postgres/txdb"
)

func TestAuthnRepository(t *testing.T) {
	db := txdb.MustOpen()
	defer db.Close()

	// migrate db
	postgres.MustMigrate(db)

	authRepo := postgres.NewAuthRepository(db)

	ctx := context.TODO()
	authKey := authn.NewAuthKey()
	err := authRepo.CreateAuthKey(ctx, authKey)
	require.NoError(t, err)
	assert.NotEmpty(t, authKey)

	exists, err := authRepo.AuthKeyExists(ctx, authKey)
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = authRepo.AuthKeyExists(ctx, authn.NewAuthKey())
	require.NoError(t, err)
	assert.False(t, exists)
}
