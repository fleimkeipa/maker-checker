package repositories

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type messageMongo struct {
	CreatedAt  time.Time          `bson:"created_at"`
	DeletedAt  time.Time          `bson:"deleted_at"`
	ID         primitive.ObjectID `bson:"_id"`
	SenderID   primitive.ObjectID `bson:"sender_id"`
	ReceiverID primitive.ObjectID `bson:"receiver_id"`
	Text       string             `bson:"text"`
	Status     int                `bson:"status"`
}
