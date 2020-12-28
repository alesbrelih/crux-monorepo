package repository

import (
	"context"
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type User struct {
	Id       int64
	Name     string
	Surname  string
	Username string
	Email    string
	Password string
}

func NewRepository(postgresDsn string) Repository {
	return &repositoryPQ{
		dsn: postgresDsn,
	}
}

type Repository interface {
	connect() *sql.DB
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	HasAccess(ctx context.Context, id int64) (bool, error)
	Migrate(path string) error
}

type repositoryPQ struct {
	dsn string
}

func (d *repositoryPQ) connect() *sql.DB {

	db, err := sql.Open("postgres", d.dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func (repo *repositoryPQ) Migrate(path string) error {
	m, err := migrate.New(
		"file://"+path,
		repo.dsn)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		return err
	}
	return nil
}

// Username can be email or username
func (repo *repositoryPQ) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	conn := repo.connect()
	defer conn.Close()

	q := `SELECT id, name, surname, username, 
				email, password
			FROM crux_user u 
			WHERE (u.email = $1 OR u.username = $1) AND u.active = true`

	var user User
	if err := conn.QueryRowContext(ctx, q, username).Scan(&user.Id, &user.Name, &user.Surname, &user.Username, &user.Email, &user.Password); err != nil {
		return nil, errors.Wrap(err, "Error retrieving user in GetUserByUsername")
	}

	return &user, nil
}

func (repo *repositoryPQ) HasAccess(ctx context.Context, id int64) (bool, error) {
	conn := repo.connect()
	defer conn.Close()

	var active bool
	q := `SELECT COUNT(*) > 0 FROM crux_user WHERE id = $1 AND active = true`

	if err := conn.QueryRowContext(ctx, q, id).Scan(&active); err != nil {
		return false, errors.Wrap(err, "Error retrieving IsUserActive")
	}
	return active, nil
}
