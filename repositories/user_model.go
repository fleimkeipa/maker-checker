package repositories

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userMongo struct {
	DeletedAt time.Time          `bson:"deleted_at"`
	CreatedAt time.Time          `bson:"created_at"`
	Username  string             `bson:"username"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	ID        primitive.ObjectID `bson:"_id"`
}
