package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/fleimkeipa/maker-checker/model"
	"github.com/fleimkeipa/maker-checker/uc"
	"github.com/fleimkeipa/maker-checker/util"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AuthHandlers struct {
	userUC *uc.UserUC
}

func NewAuthHandlers(uc *uc.UserUC) *AuthHandlers {
	return &AuthHandlers{
		userUC: uc,
	}
}

// Register godoc
//
//	@Summary		User register
//	@Description	This endpoint allows a user to log in by providing a valid username and password.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		model.Register	true	"User register input"
//	@Success		200		{object}	AuthResponse	"Successfully registered in with JWT token"
//	@Failure		400		{object}	FailureResponse	"Error message including details on failure"
//	@Failure		500		{object}	FailureResponse	"Interval error"
//	@Router			/auth/register [post]
func (rc *AuthHandlers) Register(c echo.Context) error {
	var input model.Register

	if err := c.Bind(&input); err != nil {
		var errorMessage string
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			validationError := validationErrors[0]
			if validationError.Tag() == "required" {
				errorMessage = fmt.Sprintf("%s not provided", validationError.Field())
			}
		}

		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to bind request: %v", errorMessage),
			Message: "Invalid register details. Please ensure all required fields are provided and try again.",
		})
	}

	exists, err := rc.userUC.Exists(c.Request().Context(), input.Username)
	if err != nil {
		return HandleEchoError(c, err)
	}

	if exists {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("User already exists: %v", err),
			Message: "User already exists. Please choose a different username.",
		})
	}

	if input.Password != input.ConfirmPassword {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Password does not match: %v", err),
			Message: "Password and confirm password do not match. Please try again.",
		})
	}

	newUser := model.UserCreateRequest{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}

	user, err := rc.userUC.Create(c.Request().Context(), newUser)
	if err != nil {
		return HandleEchoError(c, err)
	}

	jwt, err := util.GenerateJWT(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to generate JWT: %v", err),
			Message: "Register failed. Please try again later.",
		})
	}

	return c.JSON(http.StatusCreated, AuthResponse{
		Token:    jwt,
		Type:     "basic",
		Username: input.Username,
		Message:  "Successfully registered in",
	})
}

// Login godoc
//
//	@Summary		User login
//	@Description	This endpoint allows a user to log in by providing a valid username and password.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		model.Login		true	"User login input"
//	@Success		200		{object}	AuthResponse	"Successfully logged in with JWT token"
//	@Failure		400		{object}	FailureResponse	"Error message including details on failure"
//	@Failure		500		{object}	FailureResponse	"Interval error"
//	@Router			/auth/login [post]
func (rc *AuthHandlers) Login(c echo.Context) error {
	var input model.Login

	if err := c.Bind(&input); err != nil {
		var errorMessage string
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			validationError := validationErrors[0]
			if validationError.Tag() == "required" {
				errorMessage = fmt.Sprintf("%s not provided", validationError.Field())
			}
		}

		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to bind request: %v", errorMessage),
			Message: "Invalid login details. Please ensure all required fields are provided and try again.",
		})
	}

	user, err := rc.userUC.GetByUsernameOrEmail(c.Request().Context(), input.Username)
	if err != nil {
		return HandleEchoError(c, err)
	}

	if err := model.ValidateUserPassword(user.Password, input.Password); err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Invalid password: %v", err),
			Message: "Invalid password. Please check the password and try again.",
		})
	}

	jwt, err := util.GenerateJWT(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to generate JWT: %v", err),
			Message: "Login failed. Please try again later.",
		})
	}

	return c.JSON(http.StatusOK, AuthResponse{
		Token:    jwt,
		Type:     "basic",
		Username: input.Username,
		Message:  "Successfully logged in",
	})
}
