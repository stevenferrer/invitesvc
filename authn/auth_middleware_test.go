package authn_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stevenferrer/invitesvc/authn"
	"github.com/stevenferrer/invitesvc/postgres"
	"github.com/stevenferrer/invitesvc/postgres/txdb"
)

func TestAuthMiddleware(t *testing.T) {
	db := txdb.MustOpen()
	defer db.Close()

	// migrate db
	postgres.MustMigrate(db)

	authRepo := postgres.NewAuthRepository(db)
	authSvc := authn.NewAuthService(authRepo)

	ctx := context.Background()
	authKey, err := authSvc.GenerateAuthKey(ctx)
	require.NoError(t, err)

	authMiddleware := authn.NewAuthMiddleware(authSvc)

	// Setup
	e := echo.New()
	e.Use(authMiddleware)
	e.GET("/", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	t.Run("authorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add(authn.AuthKeyHeader, string(authKey))
		rr := httptest.NewRecorder()
		e.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

	})

	t.Run("un-authorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		e.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)

		req = httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add(authn.AuthKeyHeader, string(authn.NewAuthKey()))
		rr = httptest.NewRecorder()
		e.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

}
