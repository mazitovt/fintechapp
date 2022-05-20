package app

import (
	"github.com/mazitovt/fintechapp/auth/pkg/database/mongodb"
	"github.com/mazitovt/fintechapp/auth/repository"
	"github.com/mazitovt/fintechapp/auth/server"
	"github.com/mazitovt/fintechapp/auth/service"
	"log"
	"time"
)

func Run(config string) error {

	// TODO create index on field email: db.users.createIndex({"email": 1}, {unique: true})

	// TODO db: check tokens for expiration
	// TODO create config
	// TODO add password hashing
	// TODO parse can return parsed id or email for context propagation
	mongoUri := "mongodb://root:example@localhost:27017/"
	mongoUser := "root"
	mongoPwd := "example"
	mongoDbName := "auth_db"
	mongoCol := "users"
	accessSecret := "access_secret"
	refreshSecret := "refresh_secret"
	maxRefreshTokens := 3

	authServerAddr := "localhost:8001"
	accessTTL := 60 * time.Second
	refreshTTL := 300 * time.Second

	mongoClient, err := mongodb.NewClient(mongoUri, mongoUser, mongoPwd)
	if err != nil {
		return err
	}

	//cfg.Mongo.Name
	db := mongoClient.Database(mongoDbName)

	usersRepo := repository.NewUsersRepo(db, mongoCol)

	userService := service.NewUserService(usersRepo, maxRefreshTokens)

	tokenService := service.NewJWTokenService(accessSecret, refreshSecret)

	authServer := server.New(userService, tokenService, accessTTL, refreshTTL)

	log.Println("start server")

	if err := authServer.Run(authServerAddr); err != nil {
		log.Println("Run error: ", err)
		return err
	}

	log.Println("end server")

	return nil
}
