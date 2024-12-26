package interfaces

import (
	"context"

	"github.com/fleimkeipa/maker-checker/model"
)

type MessageInterfaces interface {
	Create(ctx context.Context, message *model.Message) (*model.Message, error)
	Update(ctx context.Context, messageID string, message *model.Message) (*model.Message, error)
	List(ctx context.Context, opts model.MessageFindOpts) ([]model.Message, error)
	GetByID(ctx context.Context, messageID string) (*model.Message, error)
}
