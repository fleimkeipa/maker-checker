package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fleimkeipa/maker-checker/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserMongoRepo struct {
	db *mongo.Database
}

func NewUserMongoRepo(db *mongo.Database) *UserMongoRepo {
	return &UserMongoRepo{
		db: db,
	}
}

var userColl = "users"

func (rc *UserMongoRepo) Create(ctx context.Context, newUser *model.User) (*model.User, error) {
	mongoUser, err := rc.internalToMongo(newUser)
	if err != nil {
		return nil, err
	}

	query, err := rc.
		db.
		Collection(userColl).
		InsertOne(ctx, mongoUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	oid, ok := query.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("can't get inserted ID")
	}

	newUser.ID = oid.Hex()

	return rc.mongoToInternal(mongoUser), nil
}

func (rc *UserMongoRepo) Update(ctx context.Context, userID string, user *model.User) (*model.User, error) {
	oID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert id: %w", err)
	}

	updateUser, err := rc.internalToMongo(user)
	if err != nil {
		return nil, fmt.Errorf("failed to convert user: %w", err)
	}

	updateDocs := bson.M{
		"$set": updateUser,
	}

	filter := bson.M{"_id": oID}
	query, err := rc.
		db.
		Collection(userColl).
		UpdateOne(ctx, updateDocs, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if query.MatchedCount == 0 {
		return nil, fmt.Errorf("not found user with id: %v", userID)
	}

	return rc.mongoToInternal(updateUser), nil
}

func (rc *UserMongoRepo) Delete(ctx context.Context, id string) error {
	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert id: %w", err)
	}

	filter := bson.M{"_id": oID}
	updater := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
			"is_active":  0,
		},
	}
	query, err := rc.
		db.
		Collection(userColl).
		UpdateOne(ctx, filter, updater)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	if query.MatchedCount == 0 {
		return fmt.Errorf("not found user with id: %v", id)
	}

	return nil
}

func (rc *UserMongoRepo) GetByID(ctx context.Context, userID string) (*model.User, error) {
	oID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert user id: %w", err)
	}

	user := new(userMongo)
	err = rc.
		db.
		Collection(userColl).
		FindOne(ctx, bson.M{"_id": oID}).
		Decode(user)
	if err != nil {
		return nil, err
	}

	return rc.mongoToInternal(user), nil
}

func (rc *UserMongoRepo) GetByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*model.User, error) {
	if usernameOrEmail == "" {
		return nil, errors.New("missing username or email")
	}

	filter := bson.M{
		"$or": []bson.M{
			{"username": usernameOrEmail},
			{"email": usernameOrEmail},
		},
	}
	user := new(userMongo)
	err := rc.
		db.
		Collection(userColl).
		FindOne(ctx, filter).
		Decode(user)
	if err != nil {
		return nil, err
	}

	return rc.mongoToInternal(user), nil
}

func (rc *UserMongoRepo) Exists(ctx context.Context, usernameOrEmail string) (bool, error) {
	if usernameOrEmail == "" {
		return false, errors.New("missing username or email")
	}

	filter := bson.M{
		"$or": []bson.M{
			{"username": usernameOrEmail},
			{"email": usernameOrEmail},
		},
	}
	count, err := rc.
		db.
		Collection(userColl).
		CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (rc *UserMongoRepo) mongoToInternal(u *userMongo) *model.User {
	return &model.User{
		CreatedAt: u.CreatedAt,
		DeletedAt: u.DeletedAt,
		ID:        u.ID.Hex(),
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		Connects:  u.Connects,
	}
}

func (rc *UserMongoRepo) internalToMongo(u *model.User) (*userMongo, error) {
	var oID primitive.ObjectID
	var err error

	if u.ID != "" {
		oID, err = primitive.ObjectIDFromHex(u.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to convert user id: %w", err)
		}
	} else {
		oID = primitive.NewObjectID()
	}

	return &userMongo{
		CreatedAt: u.CreatedAt,
		DeletedAt: u.DeletedAt,
		ID:        oID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		Connects:  u.Connects,
	}, nil
}
