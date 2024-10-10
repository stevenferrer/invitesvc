package authn

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// NewAuthKey returns an auth key handler which is used for generating auth keys
func NewAuthKeyHandler(authSvc Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		authKey, err := authSvc.GenerateAuthKey(c.Request().Context())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, echo.Map{
			"authKey": authKey,
		})
	}
}
