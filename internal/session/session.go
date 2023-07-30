package session

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/ghkadim/highload_architect/internal/models"
)

var (
	ErrTokenExpired = errors.New("token expired")
	ErrTokenInvalid = errors.New("token invalid")
)

type Session struct {
	key []byte
}

func NewSession(key string) *Session {
	return &Session{
		key: []byte(key),
	}
}

func (s *Session) HashPassword(ctx context.Context, password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (s *Session) CompareHashAndPassword(ctx context.Context, hash []byte, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *Session) TokenForUser(ctx context.Context, userID models.UserID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   string(userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	})

	return token.SignedString(s.key)
}

func (s *Session) ParseToken(ctx context.Context, tokenStr string) (models.UserID, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if token != nil && token.Valid {
		return models.UserID(claims.Subject), nil
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		return "", errors.Join(models.ErrUnauthorized, ErrTokenExpired)
	} else {
		return "", errors.Join(models.ErrUnauthorized, ErrTokenInvalid, err)
	}
}
