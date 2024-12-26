package uc

import (
	"context"
	"net/http"
	"time"

	"github.com/fleimkeipa/maker-checker/model"
	"github.com/fleimkeipa/maker-checker/pkg"
	"github.com/fleimkeipa/maker-checker/repositories/interfaces"
)

type UserUC struct {
	userRepo interfaces.UserInterfaces
}

func NewUserUC(repo interfaces.UserInterfaces) *UserUC {
	return &UserUC{
		userRepo: repo,
	}
}

func (rc *UserUC) Create(ctx context.Context, req model.UserCreateRequest) (*model.User, error) {
	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	hashedPassword, err := model.HashPassword(req.Password)
	if err != nil {
		return nil, pkg.NewError(err, "failed to hash password", http.StatusInternalServerError)
	}
	user.Password = hashedPassword

	user.CreatedAt = time.Now()

	newUser, err := rc.userRepo.Create(ctx, &user)
	if err != nil {
		return nil, pkg.NewError(err, "failed to create user", http.StatusInternalServerError)
	}

	return newUser, nil
}

func (rc *UserUC) Update(ctx context.Context, userID string, req model.UserCreateRequest) (*model.User, error) {
	// user exist control
	_, err := rc.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	hashedPassword, err := model.HashPassword(req.Password)
	if err != nil {
		return nil, pkg.NewError(err, "failed to hash password", http.StatusInternalServerError)
	}
	user.Password = hashedPassword

	updatedUser, err := rc.userRepo.Update(ctx, userID, &user)
	if err != nil {
		return nil, pkg.NewError(err, "failed to update user", http.StatusInternalServerError)
	}

	return updatedUser, nil
}

func (rc *UserUC) GetByID(ctx context.Context, id string) (*model.User, error) {
	user, err := rc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, pkg.NewError(err, "user not found", http.StatusNotFound)
	}

	return user, nil
}

func (rc *UserUC) GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*model.User, error) {
	user, err := rc.userRepo.GetByUsernameOrEmail(ctx, usernameOrEmail)
	if err != nil {
		return nil, pkg.NewError(err, "user not found", http.StatusNotFound)
	}

	return user, nil
}

func (rc *UserUC) Exists(ctx context.Context, usernameOrEmail string) (bool, error) {
	exists, err := rc.userRepo.Exists(ctx, usernameOrEmail)
	if err != nil {
		return false, pkg.NewError(err, "failed to get user by username or email", http.StatusInternalServerError)
	}

	return exists, nil
}

func (rc *UserUC) Delete(ctx context.Context, userID string) error {
	if err := rc.userRepo.Delete(ctx, userID); err != nil {
		return pkg.NewError(err, "failed to delete user", http.StatusInternalServerError)
	}

	return nil
}
