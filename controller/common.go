package controller

import (
	"strconv"

	"github.com/fleimkeipa/maker-checker/model"
	"github.com/labstack/echo/v4"
)

type FailureResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type AuthResponse struct {
	Type     string `json:"type" example:"basic,oauth2"`
	Token    string `json:"token"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func getPagination(c echo.Context) model.PaginationOpts {
	limitQuery := c.QueryParam("limit")
	skipQuery := c.QueryParam("skip")

	limit, _ := strconv.Atoi(limitQuery)

	skip, _ := strconv.Atoi(skipQuery)

	if limit <= 0 {
		limit = 30
	}

	if skip < 0 {
		skip = 0
	}

	return model.PaginationOpts{
		Skip:  uint(skip),
		Limit: uint(limit),
	}
}

func getFilter(c echo.Context, query string) model.Filter {
	param := c.QueryParam(query)
	if param == "" {
		return model.Filter{}
	}

	return model.Filter{
		IsSended: true,
		Value:    param,
	}
}
