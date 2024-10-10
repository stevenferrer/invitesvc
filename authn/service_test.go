package authn_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stevenferrer/invitesvc/authn"
	"github.com/stevenferrer/invitesvc/postgres"
	"github.com/stevenferrer/invitesvc/postgres/txdb"
)

func TestAuthService(t *testing.T) {
	db := txdb.MustOpen()
	defer db.Close()

	// migrate db
	postgres.MustMigrate(db)

	authRepo := postgres.NewAuthRepository(db)
	authSvc := authn.NewAuthService(authRepo)

	ctx := context.TODO()
	authKey, err := authSvc.GenerateAuthKey(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, authKey)

	// check auth key exists
	exists, err := authSvc.AuthKeyExists(ctx, authKey)
	require.NoError(t, err)
	assert.True(t, exists)

	// check auth key not exists
	exists, err = authRepo.AuthKeyExists(ctx, authn.NewAuthKey())
	require.NoError(t, err)
	assert.False(t, exists)
}
