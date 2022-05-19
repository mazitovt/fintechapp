package service

import (
	"context"
	"fmt"
	"github.com/mazitovt/fintechapp/auth/domain"
	"github.com/mazitovt/fintechapp/auth/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	repo             repository.Users
	maxRefreshTokens int
}

func NewUserService(repo repository.Users) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) SignUp(ctx context.Context, email, password string) (string, error) {

	passwordHash := password

	//passwordHash, err := s.hasher.Hash(password)
	//if err != nil {
	//	return nil, err
	//}

	user := domain.User{
		Password:      passwordHash,
		Email:         email,
		RefreshTokens: bson.A{},
	}

	id, err := s.repo.Create(ctx, user)
	if err != nil {
		//if errors.Is(err, domain.ErrUserAlreadyExists) {
		//	return nil, err
		//}
		return "", err
	}

	return id, nil

}

func (s *UserService) SignIn(ctx context.Context, email, password string) (string, error) {

	passwordHash := password
	//passwordHash, err := s.hasher.Hash(password)
	//if err != nil {
	//	return Tokens{}, err
	//}

	user, err := s.repo.GetByCredentials(ctx, email, passwordHash)
	if err != nil {
		//if errors.Is(err, domain.ErrUserNotFound) {
		//	return "", err
		//}

		return "", err
	}

	return user.ID.String(), nil
}

func (s *UserService) AddRefresh(ctx context.Context, userId, token string) error {
	user, err := s.repo.GetByUserId(ctx, userId)
	if err != nil {
		return fmt.Errorf("Users.GetById: %w", err)
	}

	if len(user.RefreshTokens) >= 5 {
		user.RefreshTokens = primitive.A{token}
	} else {
		user.RefreshTokens = append(user.RefreshTokens, token)
	}

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *UserService) HasToken(ctx context.Context, userId, token string) (bool, error) {
	user, err := s.repo.GetByUserId(ctx, userId)
	if err != nil {
		return false, fmt.Errorf("Users.GetById: %w", err)
	}

	for _, t := range user.RefreshTokens {
		if t == token {
			return true, nil
		}
	}

	return false, nil
}
