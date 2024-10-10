package token

import (
	"net/http"
	"time"

	"github.com/stevenferrer/invitesvc/authn"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
)

// InitAdminRoutes initializes admin routes
func InitAdminRoutes(e *echo.Echo, tokenSvc Service, authSvc authn.Service) {
	g := e.Group("/admin")
	// use auth middleware
	g.Use(authn.NewAuthMiddleware(authSvc))
	// include auth key handler for generating new authkeys
	g.POST("/authkey", authn.NewAuthKeyHandler(authSvc))

	h := &adminHandler{tokenSvc: tokenSvc}

	g.POST("/tokens", h.generateTokens)
	g.GET("/tokens", h.listTokens)
	g.GET("/tokens/:token", h.getToken)
	g.PUT("/tokens/:token/disable", h.disableToken)
}

// InitPublicRoutes initializes public routes
func InitPublicRoutes(e *echo.Echo, tokenSvc Service) {
	h := &publicHandler{tokenSvc: tokenSvc}
	e.Use(newRateLimitMiddleware())
	e.PUT("/tokens/:token/redeem", h.redeemToken)
}

// adminHandler provides admin routes
type adminHandler struct {
	tokenSvc Service
}

// genTokenResponse is the response for generating token
type genTokenResponse struct {
	Token ID `json:"token"`
}

// generateTokens handles generate token request
func (h *adminHandler) generateTokens(c echo.Context) error {
	token, err := h.tokenSvc.GenerateToken(c.Request().Context())
	if err != nil {
		return errors.Wrap(err, "generate token")
	}

	return c.JSON(http.StatusCreated, genTokenResponse{
		Token: token,
	})
}

// tokenResponse is a get token response
type tokenResponse struct {
	Token      ID        `json:"token"`
	Redeemed   bool      `json:"redeemed"`
	Disabled   bool      `json:"disabled"`
	Expiration time.Time `json:"expiration"`
}

// getToken handles get token request
func (h *adminHandler) getToken(c echo.Context) error {
	tokenID := ID(c.Param("token"))
	tk, err := h.tokenSvc.GetToken(c.Request().Context(), tokenID)
	if err != nil {
		if errors.Is(err, ErrTokenNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "token not found")
		}

		return errors.Wrap(err, "get token")
	}

	return c.JSON(http.StatusOK, tokenResponse{
		Token:      tk.ID,
		Redeemed:   tk.Redeemed(),
		Disabled:   tk.Disabled,
		Expiration: tk.Expiration(),
	})
}

// listTokens handles list token request
func (h *adminHandler) listTokens(c echo.Context) error {
	tokens, err := h.tokenSvc.ListTokens(c.Request().Context())
	if err != nil {
		return errors.Wrap(err, "list tokens")
	}

	resp := make([]tokenResponse, 0, len(tokens))
	for _, tk := range tokens {
		resp = append(resp, tokenResponse{
			Token:      tk.ID,
			Expiration: tk.Expiration(),
			Redeemed:   tk.Redeemed(),
			Disabled:   tk.Disabled,
		})
	}

	return c.JSON(http.StatusOK, resp)
}

// disableToken handles disable token request
func (h *adminHandler) disableToken(c echo.Context) error {
	tokenID := ID(c.Param("token"))
	err := h.tokenSvc.DisableToken(c.Request().Context(), tokenID)
	if err != nil {
		if errors.Is(err, ErrTokenNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "token not found")
		}

		return errors.Wrap(err, "disable token")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "token successfully disabled.",
	})
}

// publicHandler provides public routes
type publicHandler struct {
	tokenSvc Service
}

// redeemToken handles redeem token request
func (h *publicHandler) redeemToken(c echo.Context) error {
	tokenID := ID(c.Param("token"))
	err := h.tokenSvc.RedeemToken(c.Request().Context(), tokenID)
	if err != nil {
		if errors.Is(err, ErrTokenNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "token not found")
		}

		// application error
		if errors.Is(err, ErrTokenDisabled) ||
			errors.Is(err, ErrTokenExpired) ||
			errors.Is(err, ErrTokenRedeemed) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		return errors.Wrap(err, "redeem token")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "token successfully redeemed.",
	})
}

// newRateLimitMiddleware returns a new rate limit middleware
func newRateLimitMiddleware() echo.MiddlewareFunc {
	// NOTE: We can implement a custom rate-limiter store
	// instead of using the default in-memory implementation
	config := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      10,
				Burst:     30,
				ExpiresIn: 3 * time.Minute,
			},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, echo.Map{
				"message": "Too many requests",
			})
		},
	}

	return middleware.RateLimiterWithConfig(config)
}
