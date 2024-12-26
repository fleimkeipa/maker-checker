package controller

import (
	"errors"
	"net/http"

	"github.com/fleimkeipa/maker-checker/pkg"

	"github.com/labstack/echo/v4"
)

// HandleEchoError handles errors that occur within the Echo framework.
func HandleEchoError(c echo.Context, err error) error {
	var pe *pkg.Error

	if errors.As(err, &pe) {
		return c.JSON(pe.StatusCode(), FailureResponse{
			Error:   pe.Error(),
			Message: pe.Message(),
		})
	} else {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   err.Error(),
			Message: "Internal Server Error",
		})
	}
}
