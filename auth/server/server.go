package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mazitovt/fintechapp/auth/service"
	"github.com/valyala/fasthttp"
	"log"
	"time"
)

type Server struct {
	userService  service.Users
	tokenService service.Tokens
	accessTTL    time.Duration
	refreshTTL   time.Duration
}

func New(userService service.Users, tokenService service.Tokens, accessTtl, refreshTtl time.Duration) *Server {
	return &Server{
		userService:  userService,
		tokenService: tokenService,
		accessTTL:    accessTtl,
		refreshTTL:   refreshTtl,
	}
}

// TODO: check error.code and response with 409 or 500 (dont disclose internal errors)
func (s *Server) Run(addr string) error {
	if err := fasthttp.ListenAndServe(addr, func(ctx *fasthttp.RequestCtx) {

		log.Println("incoming req: ", ctx.Path())
		// if no error occurred means response body written
		code, err := s.Handle(ctx, string(ctx.Path()))
		ctx.SetStatusCode(code)

		if err != nil {
			log.Println("HandleErr: ", err)
			writeJson(ctx, map[string]any{"error": err})
		}

	}); err != nil {
		return err
	}

	return nil
}

func (s *Server) Handle(ctx *fasthttp.RequestCtx, path string) (int, error) {

	if !ctx.IsPost() {
		return fasthttp.StatusMethodNotAllowed, fmt.Errorf("wrong http method")
	}

	switch path {
	case "/signup":
		return s.signUp(ctx)
	case "/signin":
		return s.signIn(ctx)
	case "/refresh":
		return s.refresh(ctx)
	case "/parse":
		return s.parse(ctx)
	}

	return fasthttp.StatusNotFound, fmt.Errorf("no such path")
}

// doesn't have an account
// create account and return tokens or fail
func (s *Server) signUp(ctx *fasthttp.RequestCtx) (int, error) {
	args := ctx.QueryArgs()

	if !args.Has("email") || !args.Has("password") {
		return fasthttp.StatusBadRequest, fmt.Errorf("not enough arguments")
	}

	email := string(args.Peek("email"))
	password := string(args.Peek("password"))

	// add user to database
	// if email taken, write error
	userId, err := s.userService.SignUp(context.Background(), email, password)
	if err != nil {
		return fasthttp.StatusConflict, fmt.Errorf("UserService.SignUp: %w", err)
	}

	return s.writeTokens(ctx, userId)
}

// have an account
// fail or jwt pair
func (s *Server) signIn(ctx *fasthttp.RequestCtx) (int, error) {
	args := ctx.QueryArgs()

	if !args.Has("email") || !args.Has("password") {
		return fasthttp.StatusBadRequest, fmt.Errorf("not enough arguments")
	}

	email := string(args.Peek("email"))
	password := string(args.Peek("password"))

	// TODO is it right to pass Background() ???
	// check if email exists
	userId, err := s.userService.SignIn(context.Background(), email, password)
	if err != nil {
		return fasthttp.StatusBadRequest, fmt.Errorf("UserService.SignIn: %w", err)
	}

	// return user new pair of tokens
	return s.writeTokens(ctx, userId)
}

// TODO implement blacklist of refresh tokens;
// if refresh token has been stolen, i can't verify its owner
// parse refresh jwt and get new jwt pair: fail or jwt pair
func (s *Server) refresh(ctx *fasthttp.RequestCtx) (int, error) {
	// TODO where should jwt be stored: in QueryArgs or PostArgs?
	args := ctx.QueryArgs()

	if !args.Has("token") {
		return fasthttp.StatusBadRequest, fmt.Errorf("no token argument")
	}

	token := string(args.Peek("token"))

	// check if token is valid
	userId, err := s.tokenService.ParseRefresh(token)
	if err != nil {
		return fasthttp.StatusBadRequest, fmt.Errorf("Tokens.ParseRefresh: %w", err)
	}

	// check if user has token
	f, err := s.userService.HasToken(context.Background(), userId, token)
	if err != nil {
		return fasthttp.StatusInternalServerError, fmt.Errorf("Users.HasToken: %w", err)
	}

	if !f {
		return fasthttp.StatusBadRequest, fmt.Errorf("user doensn't have such token")
	}

	// return user new pair of tokens
	return s.writeTokens(ctx, userId)
}

// parse access jwt and write userID: fail or success
func (s *Server) parse(ctx *fasthttp.RequestCtx) (int, error) {
	args := ctx.QueryArgs()

	if !args.Has("token") {
		return fasthttp.StatusBadRequest, fmt.Errorf("no token argument")
	}

	token := string(args.Peek("token"))

	_, err := s.tokenService.ParseAccess(token)
	if err != nil {
		return fasthttp.StatusInternalServerError, fmt.Errorf("Tokens.ParseAccess: %w", err)
	}

	writeJson(ctx, map[string]any{"valid": true})

	return fasthttp.StatusOK, nil
}

// TODO always returns internal error
func (s *Server) writeTokens(ctx *fasthttp.RequestCtx, userId string) (int, error) {

	// TODO maybe remove ttl from function call
	// create new refresh token
	refresh, err := s.tokenService.Refresh(userId, s.refreshTTL)
	if err != nil {
		return fasthttp.StatusConflict, fmt.Errorf("TokenService.Refresh: %w", err)
	}
	// create new access token
	access, err := s.tokenService.Access(userId, s.accessTTL)
	if err != nil {
		return fasthttp.StatusConflict, fmt.Errorf("TokenService.Access: %w", err)
	}

	// add refresh token to user's token list
	// if number of token overflows limit, continue
	err = s.userService.AddRefresh(context.Background(), userId, refresh)
	if err != nil {
		return fasthttp.StatusConflict, fmt.Errorf("Users.AddRefresh: %w ", err)
	}

	writeJson(ctx, map[string]any{"access": access, "refresh": refresh})

	return fasthttp.StatusOK, nil
}

func writeJson(ctx *fasthttp.RequestCtx, value any) {
	b, _ := json.Marshal(value)
	ctx.SetContentType("application/json")
	ctx.SetBody(b)
}
