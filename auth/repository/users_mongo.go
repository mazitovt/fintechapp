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

const (
	usersCollection = "users"
)

type UsersRepo struct {
	db *mongo.Collection
}

func NewUsersRepo(db *mongo.Database) *UsersRepo {
	return &UsersRepo{
		db: db.Collection(usersCollection),
	}
}

func (r *UsersRepo) Create(ctx context.Context, user domain.User) (string, error) {
	res, err := r.db.InsertOne(ctx, user)
	if mongodb.IsDuplicate(err) {
		return "", domain.ErrUserAlreadyExists
	}

	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("ERROR: res.InsertedID.(primitive.ObjectID)")
	}

	return id.String(), err
}

func (r *UsersRepo) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	var user domain.User
	if err := r.db.FindOne(ctx, bson.M{"email": email, "password": password}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, err
	}

	return user, nil
}

func (r *UsersRepo) GetByUserId(ctx context.Context, userId string) (domain.User, error) {
	var user domain.User
	if err := r.db.FindOne(ctx, bson.D{{"_id", userId}}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, err
	}

	return user, nil
}

func (r *UsersRepo) UpdateUser(ctx context.Context, user domain.User) error {
	// TODO remove Background()
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"tokens": user.RefreshTokens}})

	return err
}
