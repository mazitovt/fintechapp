package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/mazitovt/fintechapp/auth/domain"
	"github.com/mazitovt/fintechapp/auth/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersRepo struct {
	db *mongo.Collection
}

func NewUsersRepo(db *mongo.Database, usersCollection string) *UsersRepo {
	return &UsersRepo{
		db: db.Collection(usersCollection),
	}
}

func (r *UsersRepo) Create(ctx context.Context, user domain.User) (string, error) {
	res, err := r.db.InsertOne(ctx, user)
	if mongodb.IsDuplicate(err) {
		return "", ErrUserAlreadyExists
	}

	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("ERROR: res.InsertedID.(primitive.ObjectID)")
	}

	return id.Hex(), err
}

func (r *UsersRepo) GetByCredentials(ctx context.Context, email, password string) (user domain.User, err error) {
	if err = r.db.FindOne(ctx, bson.M{"email": email, "password": password}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = ErrUserNotFound
		}

		return
	}

	return
}

func (r *UsersRepo) GetByUserId(ctx context.Context, userId string) (user domain.User, err error) {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return
	}
	if err = r.db.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = ErrUserNotFound
		}

		return
	}

	return
}

func (r *UsersRepo) UpdateUser(ctx context.Context, user domain.User) error {
	// TODO maybe append to array, and not rewrite it
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"tokens": user.RefreshTokens}})

	return err
}
