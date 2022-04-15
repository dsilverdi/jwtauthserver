package database

import (
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
)

func Connect() (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error
	url := "root" + ":" + "jwtauthserver" + "@tcp(" + "jwtauthdb" + ":" + "3306" + ")/" + "jwtauth" + "?parseTime=true&clientFoundRows=true"

	for {
		db, err = sqlx.Connect("mysql", url)
		if err == nil {
			break
		}

		if !strings.Contains(err.Error(), "connect: connection refused") {
			return nil, err
		}

		const retryDuration = 5 * time.Second
		time.Sleep(retryDuration)
	}

	if err := migrateDB(db); err != nil {
		return nil, err
	}

	return db, nil
}

func migrateDB(db *sqlx.DB) error {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "users_1",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS users (
						id       	VARCHAR(255),
						username    VARCHAR(255) NOT NULL UNIQUE,
						password	VARCHAR(255) NOT NULL,
						created_at  DATETIME,
						PRIMARY KEY (id)
					) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;;`,

					// `CREATE TABLE IF NOT EXISTS users_auth (
					// 	id       		VARCHAR(255),
					// 	access_token    VARCHAR(255) NOT NULL,
					// 	expired		    INT NOT NULL,
					// 	created_at		DATETIME,
					// 	PRIMARY KEY (id)
					// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;;`,
				},
				Down: []string{
					"DROP TABLE users",
				},
			},
		},
	}

	_, err := migrate.Exec(db.DB, "mysql", migrations, migrate.Up)
	return err
}
