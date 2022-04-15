package auth

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	ID        string
	Username  string
	Password  string
	CreatedAt time.Time
}

type Auth struct {
	Token    string
	Username string
}

type JwtClaims struct {
	jwt.StandardClaims
	Username string
}

type UserRepository interface {
	Save(ctx context.Context, user User) error
	Read(ctx context.Context, username string) (*User, error)
}
