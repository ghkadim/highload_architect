package session

import (
	"context"
	"errors"
	"github.com/ghkadim/highload_architect/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
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

func (s *Session) TokenForUser(ctx context.Context, user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"nbf":    time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(s.key)
}
