package app

import (
	"github.com/mazitovt/fintechapp/auth/pkg/database/mongodb"
	"github.com/mazitovt/fintechapp/auth/repository"
	"github.com/mazitovt/fintechapp/auth/server"
	"github.com/mazitovt/fintechapp/auth/service"
	"log"
)

func Run(config string) error {

	//cfg.Mongo.URI, cfg.Mongo.User, cfg.Mongo.Password
	mongoClient, err := mongodb.NewClient("mongodb://root:example@mongo:27017/", "root", "example")
	if err != nil {
		return err
	}

	//cfg.Mongo.Name
	db := mongoClient.Database("auth_db")

	usersRepo := repository.NewUsersRepo(db)

	userService := service.NewUserService(usersRepo)

	tokenService := service.NewJWTokenService("access_secret", "refresh_secret")

	authServer := server.New(userService, tokenService)

	if err := authServer.Run("localhost:8001"); err != nil {
		log.Println("Run error: ", err)
		return err
	}

	return nil
}
