package service

import (
	"context"
	"time"
)

type JWTTokens struct {
	accessToken  string
	refreshToken string
}

type Users interface {

	// SignUp creates new user and return userId
	SignUp(ctx context.Context, email, password string) (string, error)

	// SignIn return userId or err if email isn't found in db
	SignIn(ctx context.Context, email, password string) (string, error)

	// AddRefresh throws error when number of refresh tokens overflows limit
	AddRefresh(ctx context.Context, userId, token string) error

	// HasToken tells if userId has this token
	HasToken(ctx context.Context, userId, token string) (bool, error)
	//Verify(ctx context.Context, userID primitive.ObjectID, hash string) error
}

type Tokens interface {
	Access(sub string, ttl time.Duration) (string, error)
	Refresh(sub string, ttl time.Duration) (string, error)
	ParseAccess(token string) (string, error)
	ParseRefresh(token string) (string, error)
}
