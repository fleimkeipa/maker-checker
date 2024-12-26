package util

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fleimkeipa/maker-checker/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// retrieve JWT key from .env file
var privateKey = []byte("secret")

// GenerateJWT generate JWT token
func GenerateJWT(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"iat":      time.Now().Unix(),
		"eat":      time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString(privateKey)
}

// validate JWT token
func ValidateJWT(c echo.Context) error {
	token, err := getToken(c)
	if err != nil {
		return err
	}

	if token == nil || !token.Valid {
		return errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	if claims["eat"].(float64) < float64(time.Now().Unix()) {
		return errors.New("token expired")
	}

	if claims == nil {
		return errors.New("invalid token claims. claims is nil")
	}

	return nil
}

// GetOwnerFromToken returns the owner details from the JWT token
func GetOwnerFromToken(c echo.Context) (model.TokenOwner, error) {
	token, err := getToken(c)
	if err != nil {
		return model.TokenOwner{}, err
	}

	if !token.Valid {
		return model.TokenOwner{}, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return model.TokenOwner{}, errors.New("invalid token claims")
	}

	id, ok := claims["id"].(string)
	if !ok {
		return model.TokenOwner{}, errors.New("invalid id claims")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return model.TokenOwner{}, errors.New("invalid username claims")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return model.TokenOwner{}, errors.New("invalid email claims")
	}

	return model.TokenOwner{
		ID:       id,
		Username: username,
		Email:    email,
	}, nil
}

// GetOwnerIDFromCtx returns the owner id from the context string type
func GetOwnerIDFromCtx(ctx context.Context) string {
	owner, ok := ctx.Value("user").(model.TokenOwner)
	if ok {
		return owner.ID
	}

	return ""
}

// check token validity
func getToken(context echo.Context) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(context)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return privateKey, nil
	})

	return token, err
}

// extract token from request Authorization header
func getTokenFromRequest(c echo.Context) string {
	bearerToken := c.Request().Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}
