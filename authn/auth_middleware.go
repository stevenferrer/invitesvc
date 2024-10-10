package authn

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// AuthKeyHeader is the auth key header
const AuthKeyHeader = "X-AUTH-KEY"

// NewAuthMiddleware is an authentication middleware
func NewAuthMiddleware(authSvc Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authKey := c.Request().Header.Get(AuthKeyHeader)
			if authKey == "" {
				return echo.ErrUnauthorized
			}

			// check auth key exists
			exists, err := authSvc.AuthKeyExists(c.Request().Context(), AuthKey(authKey))
			if err != nil {
				return errors.Wrap(err, "check auth key exists")
			}

			if !exists {
				return echo.ErrUnauthorized
			}

			return next(c)
		}
	}
}
