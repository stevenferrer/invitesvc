package inmem_test

import (
	"context"
	"testing"

	"github.com/hashicorp/go-memdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stevenferrer/invitesvc/authn"
	"github.com/stevenferrer/invitesvc/inmem"
)

func TestAuthnRepository(t *testing.T) {
	db, err := memdb.NewMemDB(inmem.Schema())
	require.NoError(t, err)

	authRepo := inmem.NewAuthRepository(db)

	ctx := context.TODO()
	authKey := authn.NewAuthKey()
	err = authRepo.CreateAuthKey(ctx, authKey)
	require.NoError(t, err)
	assert.NotEmpty(t, authKey)

	exists, err := authRepo.AuthKeyExists(ctx, authKey)
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = authRepo.AuthKeyExists(ctx, authn.NewAuthKey())
	require.NoError(t, err)
	assert.False(t, exists)
}
