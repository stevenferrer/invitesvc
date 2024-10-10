package token_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stevenferrer/invitesvc/token"
)

func TestToken(t *testing.T) {
	tokenID, err := token.NewID()
	require.NoError(t, err)
	assert.NotEmpty(t, tokenID)

	createdAt := time.Date(2021, time.August, 13, 11, 0, 0, 0, time.Local)
	token := &token.Token{
		ID:         tokenID,
		CreatedAt:  &createdAt,
		RedeemedAt: &createdAt,
	}

	expiration := createdAt.Add(time.Hour * 24 * 7)
	assert.Equal(t, token.Expiration(), expiration)
	assert.True(t, token.Expired())
	assert.True(t, token.Redeemed())
	assert.Error(t, token.Validate())
}
