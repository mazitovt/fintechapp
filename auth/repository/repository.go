package repository

import (
	"context"
	"github.com/mazitovt/fintechapp/auth/domain"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Users interface {
	Create(ctx context.Context, user domain.User) (string, error)
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
	GetByUserId(ctx context.Context, userId string) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) error

	//Verify(ctx context.Context, userID primitive.ObjectID, code string) error
}
