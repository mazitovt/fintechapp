package service

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"math/rand"
	"time"
)

type (
	ErrEmptyToken      struct{}
	ErrEmptySigningKey struct{}
	ErrInvalidToken    struct{}
	ErrInternal        struct{}
)

type JWTokenService struct {
	accessSigningKey  string
	refreshSigningKey string
}

func NewJWTokenService(accessSigningKey, refreshSigningKey string) *JWTokenService {
	return &JWTokenService{accessSigningKey: accessSigningKey, refreshSigningKey: refreshSigningKey}
}

func (m *JWTokenService) Access(sub string, ttl time.Duration) (string, error) {
	return m.token(sub, ttl, m.accessSigningKey)
}

func (m JWTokenService) Refresh(sub string, ttl time.Duration) (string, error) {
	return m.token(sub, ttl, m.refreshSigningKey)
}

func (m *JWTokenService) ParseAccess(tkn string) (string, error) {
	return m.parse(tkn, m.accessSigningKey)
}

func (m *JWTokenService) ParseRefresh(tkn string) (string, error) {
	return m.parse(tkn, m.refreshSigningKey)
}

func (m *JWTokenService) parse(tkn string, signingKey string) (string, error) {

	if tkn == "" {
		return "", fmt.Errorf("empty access token")
	}
	if signingKey == "" {
		return "", fmt.Errorf("empty signingKey")
	}

	var claims jwt.RegisteredClaims

	token, err :=
		jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()})).
			ParseWithClaims(tkn, &claims, func(token *jwt.Token) (interface{}, error) { return []byte(signingKey), nil })

	if err != nil || !token.Valid {
		return "", err
	}

	return claims.Subject, nil
}

func (JWTokenService) token(sub string, ttl time.Duration, signingKey string) (string, error) {

	if signingKey == "" {
		return "", fmt.Errorf("empty signingKey")
	}

	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   sub,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl))},
	)

	return tkn.SignedString([]byte(signingKey))
}

func (m JWTokenService) refreshOld() string {
	p := make([]byte, 32)
	rand.New(rand.NewSource(time.Now().Unix())).Read(p)
	return fmt.Sprintf("%x", p)
}
