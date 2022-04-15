package auth

import (
	"context"
	"crypto/sha256"
	stderr "errors"
	"fmt"
	"jwtauthserver"
	"jwtauthserver/pkg/errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-sql-driver/mysql"
)

var JWT_SECRET_KEY = []byte("jwtauthserver-signature-key")

type Service interface {
	Authorize(ctx context.Context, username, password string) (*Auth, error)
	Register(ctx context.Context, username, password string) error
	IdentifyUser(ctx context.Context, userid string) (*User, error)
	TokenValidation(ctx context.Context, token string) (string, error)
}

type AuthService struct {
	User       UserRepository
	IDProvider jwtauthserver.IDprovider
}

func NewService(user UserRepository, idprov jwtauthserver.IDprovider) Service {
	return &AuthService{
		User:       user,
		IDProvider: idprov,
	}
}

func (svc *AuthService) Authorize(ctx context.Context, username string, password string) (*Auth, error) {
	CurrentUser, err := svc.User.Read(ctx, username)
	if err != nil {
		return nil, errors.Wrap(errors.ErrNotFound, err)
	}

	userpwd := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	if userpwd != CurrentUser.Password {
		return nil, errors.ErrWrongPassword
	}

	claims := &JwtClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "JWT_AUTH_SERVER",
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		},
		Username: CurrentUser.Username,
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString(JWT_SECRET_KEY)
	if err != nil {
		return nil, err
	}

	return &Auth{
		Token: tokenString,
	}, nil
}

func (svc *AuthService) Register(ctx context.Context, username string, password string) error {
	var mysqlErr *mysql.MySQLError

	NewUser := &User{
		Username:  username,
		Password:  fmt.Sprintf("%x", sha256.Sum256([]byte(password))),
		CreatedAt: time.Now(),
	}

	id, err := svc.IDProvider.ID()
	if err != nil {
		return errors.Wrap(errors.ErrCreateUUID, err)
	}

	NewUser.ID = id

	// Perform DB Call Here
	err = svc.User.Save(ctx, *NewUser)
	if err != nil {
		if stderr.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return errors.Wrap(errors.ErrAlreadyExists, err)
		}
		return errors.Wrap(errors.ErrCreateEntity, err)
	}

	return nil
}

func (svc *AuthService) IdentifyUser(ctx context.Context, ID string) (*User, error) {
	var User User
	return &User, nil
}

func (svc *AuthService) TokenValidation(ctx context.Context, token string) (string, error) {
	return "&User", nil
}
