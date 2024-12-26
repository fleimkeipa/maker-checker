package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/fleimkeipa/maker-checker/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MsgMongoRepo struct {
	db *mongo.Database
}

func NewMsgMongoRepo(db *mongo.Database) *MsgMongoRepo {
	return &MsgMongoRepo{
		db: db,
	}
}

var msgColl = "messages"

func (rc *MsgMongoRepo) Create(ctx context.Context, newMessage *model.Message) (*model.Message, error) {
	mongoMsg, err := rc.internalToMongo(newMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to convert message: %w", err)
	}

	query, err := rc.
		db.
		Collection(msgColl).
		InsertOne(ctx, &mongoMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	oid, ok := query.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("can't get inserted ID")
	}

	newMessage.ID = oid.Hex()

	return newMessage, nil
}

func (rc *MsgMongoRepo) Update(ctx context.Context, msgID string, message *model.Message) (*model.Message, error) {
	oID, err := primitive.ObjectIDFromHex(msgID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert message id: %w", err)
	}

	filter := bson.M{"_id": oID}
	update := bson.M{
		"$set": bson.M{
			"status": message.Status,
		},
	}
	query, err := rc.
		db.
		Collection(msgColl).
		UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	if query.MatchedCount == 0 {
		return nil, fmt.Errorf("not found message with id: %v", msgID)
	}

	return message, nil
}

func (rc *MsgMongoRepo) List(ctx context.Context, opts model.MessageFindOpts) ([]model.Message, error) {
	filter := bson.M{}
	if opts.ReceiverID.IsSended {
		oID, err := primitive.ObjectIDFromHex(opts.ReceiverID.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to convert receiver id: %w", err)
		}
		filter["receiver_id"] = oID
		filter["status"] = model.MessageStatusAccepted
	}
	if opts.SenderID.IsSended {
		oID, err := primitive.ObjectIDFromHex(opts.SenderID.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to convert sender id: %w", err)
		}
		filter["sender_id"] = oID
	}
	if opts.Status.IsSended {
		filter["status"] = opts.Status.Value
	}

	mongoOptions := options.Find().
		SetLimit(int64(opts.Limit)).
		SetSkip(int64(opts.Skip))

	msgs := make([]messageMongo, 0)
	cur, err := rc.
		db.
		Collection(msgColl).
		Find(ctx, filter, mongoOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find messages: %w", err)
	}

	if err := cur.All(ctx, &msgs); err != nil {
		return nil, fmt.Errorf("failed to decode messages: %w", err)
	}

	res := make([]model.Message, 0, len(msgs))
	for _, v := range msgs {
		res = append(res, *rc.mongoToInternal(&v))
	}

	return res, nil
}

func (rc *MsgMongoRepo) GetByID(ctx context.Context, msgID string) (*model.Message, error) {
	oID, err := primitive.ObjectIDFromHex(msgID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert message id: %w", err)
	}

	msg := new(messageMongo)
	err = rc.
		db.
		Collection(msgColl).
		FindOne(ctx, bson.M{"_id": oID}).
		Decode(msg)
	if err != nil {
		return nil, err
	}

	return rc.mongoToInternal(msg), nil
}

func (rc *MsgMongoRepo) mongoToInternal(msg *messageMongo) *model.Message {
	return &model.Message{
		CreatedAt:  msg.CreatedAt,
		DeletedAt:  msg.DeletedAt,
		ID:         msg.ID.Hex(),
		SenderID:   msg.SenderID.Hex(),
		ReceiverID: msg.ReceiverID.Hex(),
		Text:       msg.Text,
		Status:     msg.Status,
	}
}

func (rc *MsgMongoRepo) internalToMongo(msg *model.Message) (*messageMongo, error) {
	var mID primitive.ObjectID
	var err error

	if msg.ID != "" {
		mID, err = primitive.ObjectIDFromHex(msg.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to convert message id: %w", err)
		}
	} else {
		mID = primitive.NewObjectID()
	}

	senderID, err := primitive.ObjectIDFromHex(msg.SenderID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert sender id: %w", err)
	}
	receiverID, err := primitive.ObjectIDFromHex(msg.ReceiverID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert receiver id: %w", err)
	}

	return &messageMongo{
		CreatedAt:  msg.CreatedAt,
		DeletedAt:  msg.DeletedAt,
		ID:         mID,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Text:       msg.Text,
		Status:     msg.Status,
	}, nil
}
