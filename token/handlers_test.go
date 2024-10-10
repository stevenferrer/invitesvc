package token_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stevenferrer/invitesvc/authn"
	"github.com/stevenferrer/invitesvc/postgres"
	"github.com/stevenferrer/invitesvc/postgres/txdb"
	"github.com/stevenferrer/invitesvc/token"
)

func TestHandlers(t *testing.T) {
	db := txdb.MustOpen()
	defer db.Close()

	// migrate db
	postgres.MustMigrate(db)

	tokenRepo := postgres.NewTokenRepository(db)
	tokenSvc := token.NewService(tokenRepo)

	authRepo := postgres.NewAuthRepository(db)
	authSvc := authn.NewAuthService(authRepo)

	ctx := context.TODO()
	authKey, err := authSvc.GenerateAuthKey(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, authKey)

	e := echo.New()
	token.InitAdminRoutes(e, tokenSvc, authSvc)
	token.InitPublicRoutes(e, tokenSvc)

	t.Run("generate and retrieve token", func(t *testing.T) {
		// generate token
		req := httptest.NewRequest(http.MethodPost, "/admin/tokens", nil)
		req.Header.Add(authn.AuthKeyHeader, string(authKey))

		rr := httptest.NewRecorder()

		e.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)

		var resp1 = struct {
			Token string `json:"token"`
		}{}
		err = json.NewDecoder(rr.Body).Decode(&resp1)
		require.NoError(t, err)
		assert.NotEmpty(t, resp1.Token)

		// retrieve token
		req = httptest.NewRequest(http.MethodGet, "/admin/tokens/"+resp1.Token, nil)
		req.Header.Add(authn.AuthKeyHeader, string(authKey))
		rr = httptest.NewRecorder()

		e.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var resp2 = struct {
			Token      string    `json:"token"`
			Disabled   bool      `json:"disabled"`
			Redeemed   bool      `json:"redeemed"`
			Expiration time.Time `json:"expiration"`
		}{}
		err = json.NewDecoder(rr.Body).Decode(&resp2)
		require.NoError(t, err)

		assert.Equal(t, resp1.Token, resp2.Token)
		assert.NotZero(t, resp2.Expiration)
		assert.False(t, resp2.Redeemed)
		assert.False(t, resp2.Disabled)

		t.Run("token not found", func(t *testing.T) {
			tokenID, err := token.NewID()
			require.NoError(t, err)

			req = httptest.NewRequest(http.MethodGet, "/admin/tokens/"+string(tokenID), nil)
			req.Header.Add(authn.AuthKeyHeader, string(authKey))
			rr = httptest.NewRecorder()

			e.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusNotFound, rr.Code)
		})
	})

	t.Run("list tokens", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/tokens", nil)
		req.Header.Add(authn.AuthKeyHeader, string(authKey))

		rr := httptest.NewRecorder()

		e.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var resp = []struct {
			Token      string    `json:"token"`
			Redeemed   bool      `json:"redeemed"`
			Disabled   bool      `json:"disabled"`
			Expiration time.Time `json:"expiration"`
		}{}
		err = json.NewDecoder(rr.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Len(t, resp, 1)

		for _, res := range resp {
			assert.NotEmpty(t, res.Token)
			assert.False(t, res.Redeemed)
			assert.False(t, res.Disabled)
			assert.NotZero(t, res.Expiration)
		}
	})

	t.Run("disable token", func(t *testing.T) {
		tokens, err := tokenSvc.ListTokens(ctx)
		require.NoError(t, err)
		assert.Len(t, tokens, 1)
		// use one token for testing
		tk := tokens[0]
		urlStr := fmt.Sprintf("/admin/tokens/%s/disable", tk.ID)
		req := httptest.NewRequest(http.MethodPut, urlStr, nil)
		req.Header.Add(authn.AuthKeyHeader, string(authKey))

		rr := httptest.NewRecorder()

		e.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		// verify token update
		gotTk, err := tokenSvc.GetToken(ctx, tk.ID)
		require.NoError(t, err)
		assert.True(t, gotTk.Disabled)

		t.Run("token not found", func(t *testing.T) {
			tokenID, err := token.NewID()
			require.NoError(t, err)

			urlStr = fmt.Sprintf("/admin/tokens/%s/disable", tokenID)
			req = httptest.NewRequest(http.MethodPut, urlStr, nil)
			req.Header.Add(authn.AuthKeyHeader, string(authKey))
			rr = httptest.NewRecorder()

			e.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusNotFound, rr.Code)
		})
	})

	t.Run("redeem token", func(t *testing.T) {
		tk1, err := tokenSvc.GenerateToken(ctx)
		require.NoError(t, err)

		urlStr := fmt.Sprintf("/tokens/%s/redeem", tk1)
		req := httptest.NewRequest(http.MethodPut, urlStr, nil)
		rr := httptest.NewRecorder()

		e.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		t.Run("error redeem", func(t *testing.T) {
			urlStr = fmt.Sprintf("/tokens/%s/redeem", tk1)
			req = httptest.NewRequest(http.MethodPut, urlStr, nil)
			rr = httptest.NewRecorder()

			e.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		})

		t.Run("rate limit", func(t *testing.T) {
			for i := 0; i < 25; i++ {
				urlStr = fmt.Sprintf("/tokens/%s/redeem", tk1)
				req = httptest.NewRequest(http.MethodPut, urlStr, nil)
				rr = httptest.NewRecorder()

				e.ServeHTTP(rr, req)
			}

			urlStr = fmt.Sprintf("/tokens/%s/redeem", tk1)
			req = httptest.NewRequest(http.MethodPut, urlStr, nil)
			rr = httptest.NewRecorder()

			e.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusTooManyRequests, rr.Code)
		})
	})
}
