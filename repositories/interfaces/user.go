package interfaces

import (
	"context"

	"github.com/fleimkeipa/maker-checker/model"
)

type UserInterfaces interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Update(ctx context.Context, userID string, user *model.User) (*model.User, error)
	GetByID(ctx context.Context, userID string) (*model.User, error)
	GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*model.User, error)
	Exists(ctx context.Context, usernameOrEmail string) (bool, error)
	Delete(ctx context.Context, userID string) error
}
