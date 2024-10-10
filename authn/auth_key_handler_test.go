package authn_test

import (
	"encoding/json"
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

func TestAuthKeyHandler(t *testing.T) {
	db := txdb.MustOpen()
	defer db.Close()

	// migrate db
	postgres.MustMigrate(db)

	authRepo := postgres.NewAuthRepository(db)
	authSvc := authn.NewAuthService(authRepo)

	authKeyHandler := authn.NewAuthKeyHandler(authSvc)
	e := echo.New()
	e.POST("/", authKeyHandler)

	// generate token
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rr := httptest.NewRecorder()
	e.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)

	var resp = struct {
		AuthKey string `json:"authKey"`
	}{}
	err := json.NewDecoder(rr.Body).Decode(&resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.AuthKey)
}
