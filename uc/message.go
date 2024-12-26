package uc

import (
	"context"
	"net/http"
	"time"

	"github.com/fleimkeipa/maker-checker/model"
	"github.com/fleimkeipa/maker-checker/pkg"
	"github.com/fleimkeipa/maker-checker/repositories/interfaces"
	"github.com/fleimkeipa/maker-checker/util"
)

type MsgUC struct {
	msgRepo interfaces.MessageInterfaces
}

func NewMessageUC(repo interfaces.MessageInterfaces) *MsgUC {
	return &MsgUC{
		msgRepo: repo,
	}
}

func (rc *MsgUC) Create(ctx context.Context, req *model.MessageCreateRequest) (*model.Message, error) {
	message := model.Message{
		CreatedAt:  time.Now(),
		SenderID:   util.GetOwnerIDFromCtx(ctx),
		ReceiverID: req.ReceiverID,
		Text:       req.Text,
		Status:     model.MessageStatusPending,
	}

	newMsg, err := rc.msgRepo.Create(ctx, &message)
	if err != nil {
		return nil, pkg.NewError(err, "failed to create message", http.StatusInternalServerError)
	}

	return newMsg, nil
}

func (rc *MsgUC) Update(ctx context.Context, messageID string, req *model.MessageUpdateRequest) (*model.Message, error) {
	// message exist control
	message, err := rc.GetByID(ctx, messageID)
	if err != nil {
		return nil, err
	}

	if message.Status != model.MessageStatusPending {
		return nil, pkg.NewError(nil, "message is not pending", http.StatusConflict)
	}

	message.Status = req.Status

	_, err = rc.msgRepo.Update(ctx, messageID, message)
	if err != nil {
		return nil, pkg.NewError(err, "failed to update message", http.StatusInternalServerError)
	}

	return message, nil
}

func (rc *MsgUC) List(ctx context.Context, opts model.MessageFindOpts) ([]model.Message, error) {
	message, err := rc.msgRepo.List(ctx, opts)
	if err != nil {
		return nil, pkg.NewError(err, "messages not found", http.StatusNotFound)
	}

	return message, nil
}

func (rc *MsgUC) GetByID(ctx context.Context, messageID string) (*model.Message, error) {
	message, err := rc.msgRepo.GetByID(ctx, messageID)
	if err != nil {
		return nil, pkg.NewError(err, "message not found", http.StatusNotFound)
	}

	return message, nil
}
