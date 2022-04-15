package database

import (
	"context"
	"jwtauthserver/auth"
	"time"
)

type UserRepository struct {
	db Database
}
type UserDB struct {
	ID        string    `db:"id"`
	UserName  string    `db:"username"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

func NewUsersRepository(db Database) auth.UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) Save(ctx context.Context, user auth.User) error {
	query := `INSERT INTO users (id, username, password, created_at)
	VALUES (:id, :username, :password, :created_at);`

	userDB := &UserDB{
		ID:        user.ID,
		UserName:  user.Username,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
	}

	_, err := u.db.NamedExecContext(ctx, query, userDB)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) Read(ctx context.Context, username string) (*auth.User, error) {
	var userDB UserDB

	query := `SELECT id, username, password, created_at FROM users WHERE username = ?`

	err := u.db.QueryRowxContext(ctx, query, username).StructScan(&userDB)
	if err != nil {
		return nil, err
	}

	User := &auth.User{
		ID:        userDB.ID,
		Username:  userDB.UserName,
		Password:  userDB.Password,
		CreatedAt: userDB.CreatedAt,
	}

	return User, nil
}
