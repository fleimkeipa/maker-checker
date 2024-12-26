package util

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

// check for valid user token
func JWTAuthUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := ValidateJWT(c); err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "Authentication required",
				"error":   err.Error(),
			})
		}

		setOwnerOnCtx(c)

		return next(c)
	}
}

func setOwnerOnCtx(c echo.Context) error {
	user, err := GetOwnerFromToken(c)
	if err != nil {
		return err
	}

	ctx := context.WithValue(c.Request().Context(), "user", user)

	c.SetRequest(c.Request().WithContext(ctx))
	return nil
}
