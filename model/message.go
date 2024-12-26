package model

import "time"

const (
	MessageStatusPending  = 1
	MessageStatusAccepted = 2
	MessageStatusRejected = 3
)

type Message struct {
	CreatedAt  time.Time `json:"created_at"`
	DeletedAt  time.Time `json:"deleted_at"`
	ID         string    `json:"id"`
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Text       string    `json:"text"`
	Status     int       `json:"status"`
}

type MessageCreateRequest struct {
	ReceiverID string `json:"receiver_id"`
	Text       string `json:"text"`
}

type MessageUpdateRequest struct {
	Status int `json:"status"`
}

type MessageFindOpts struct {
	PaginationOpts
	ReceiverID Filter
	SenderID   Filter
	Status     Filter
}
