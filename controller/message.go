package controller

import (
	"fmt"
	"net/http"

	"github.com/fleimkeipa/maker-checker/model"
	"github.com/fleimkeipa/maker-checker/uc"

	"github.com/labstack/echo/v4"
)

type MessageHandlers struct {
	msgUC *uc.MsgUC
}

func NewMessageHandlers(uc *uc.MsgUC) *MessageHandlers {
	return &MessageHandlers{
		msgUC: uc,
	}
}

// Create godoc
//
//	@Summary		Create creates a new message
//	@Description	This endpoint creates a new message by providing sender id, receiver id, text, and status.
//	@Tags			messages
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			body	body		model.MessageCreateRequest	true	"Message creation input"
//	@Success		201		{object}	SuccessResponse				"message id"
//	@Failure		400		{object}	FailureResponse				"Error message including details on failure"
//	@Failure		500		{object}	FailureResponse				"Interval error"
//	@Router			/messages [post]
func (rc *MessageHandlers) Create(c echo.Context) error {
	req := new(model.MessageCreateRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to bind request: %v", err),
			Message: "Invalid request data. Please check your input and try again.",
		})
	}

	newMessage, err := rc.msgUC.Create(c.Request().Context(), req)
	if err != nil {
		return HandleEchoError(c, err)
	}

	return c.JSON(http.StatusCreated, SuccessResponse{
		Data:    newMessage.ID,
		Message: "Message created successfully.",
	})
}

// Update godoc
//
//	@Summary		Update updates an existing message
//	@Description	This endpoint updates a message by providing message id, sender id, receiver id, text, and status.
//	@Tags			messages
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		string						true	"Message id"
//	@Param			body	body		model.MessageUpdateRequest	true	"Message update input, status= pending:1, approved:2, rejected:3"
//	@Success		200		{object}	SuccessResponse				"message id"
//	@Failure		400		{object}	FailureResponse				"Error message including details on failure"
//	@Failure		500		{object}	FailureResponse				"Interval error"
//	@Router			/messages/{id} [patch]
func (rc *MessageHandlers) Update(c echo.Context) error {
	id := c.Param("id")
	input := new(model.MessageUpdateRequest)

	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to bind request: %v", err),
			Message: "Invalid request data. Please check your input and try again.",
		})
	}

	message, err := rc.msgUC.Update(c.Request().Context(), id, input)
	if err != nil {
		return HandleEchoError(c, err)
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    message.ID,
		Message: "Message updated successfully.",
	})
}

// List godoc
//
//	@Summary		List lists messages
//	@Description	This endpoint lists messages by providing limit and skip.
//	@Tags			messages
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			limit		query		int				false	"Messages limit"
//	@Param			skip		query		int				false	"Skip messages"
//	@Param			receiver_id	query		string			false	"Receiver id"
//	@Param			sender_id	query		string			false	"Sender id"
//	@Param			status		query		string			false	"Status"
//	@Success		200			{object}	SuccessResponse	"messages"
//	@Failure		400			{object}	FailureResponse	"Error message including details on failure"
//	@Failure		500			{object}	FailureResponse	"Interval error"
//	@Router			/messages [get]
func (rc *MessageHandlers) List(c echo.Context) error {
	opts := getMessageFindOpts(c)

	message, err := rc.msgUC.List(c.Request().Context(), opts)
	if err != nil {
		return HandleEchoError(c, err)
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    message,
		Message: "Message retrieved successfully.",
	})
}

// GetByID godoc
//
//	@Summary		GetByID gets a message by id
//	@Description	This endpoint gets a message by providing message id.
//	@Tags			messages
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string			true	"Message id"
//	@Success		200	{object}	SuccessResponse	"message"
//	@Failure		400	{object}	FailureResponse	"Error message including details on failure"
//	@Failure		500	{object}	FailureResponse	"Interval error"
//	@Router			/messages/{id} [get]
func (rc *MessageHandlers) GetByID(c echo.Context) error {
	id := c.Param("id")

	message, err := rc.msgUC.GetByID(c.Request().Context(), id)
	if err != nil {
		return HandleEchoError(c, err)
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    message,
		Message: "Message retrieved successfully.",
	})
}

func getMessageFindOpts(c echo.Context) model.MessageFindOpts {
	return model.MessageFindOpts{
		PaginationOpts: getPagination(c),
		ReceiverID:     getFilter(c, "receiver_id"),
		SenderID:       getFilter(c, "sender_id"),
		Status:         getFilter(c, "status"),
	}
}
