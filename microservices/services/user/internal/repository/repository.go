package repository

import (
	"context"

	"github.com/alesbrelih/crux-monorepo/microservices/internal/repository_connect"
	"github.com/alesbrelih/crux-monorepo/microservices/services/user/internal/models"
	"github.com/pkg/errors"
)

func NewRepository(dsn string) Repository {
	return &repositoryPQ{
		conn: repository_connect.NewRepositoryConnect(dsn),
	}
}

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) (string, error)
	DeleteUser(ctx context.Context, id string) error
}

type repositoryPQ struct {
	conn repository_connect.RepositoryConnect
}

func (repo *repositoryPQ) CreateUser(ctx context.Context, user *models.User) (string, error) {
	conn, err := repo.conn.Connect()
	defer conn.Close()

	if err != nil {
		return "", errors.Wrap(err, "CreateUser - Db connecting error")
	}

	q := `INSERT INTO crux_user (id, email, username, pass)
			VALUES ($1, $2, $3, $4)`
	_, err = conn.ExecContext(ctx, q, user.Id, user.Email, user.Username, user.Password)

	if err != nil {
		return "", errors.Wrap(err, "CreateUser - Inserting user error")
	}

	// this isnt needed but if implementation changes to user autoincrement
	// its better to have it defined this way already
	return user.Id, nil
}

func (repo *repositoryPQ) DeleteUser(ctx context.Context, id string) error {
	conn, err := repo.conn.Connect()
	defer conn.Close()

	if err != nil {
		return errors.Wrap(err, "DeleteUser - Db connecting error")
	}

	q := "DELETE FROM crux_user WHERE id := $1"
	_, err = conn.ExecContext(ctx, q, id)

	if err != nil {
		return errors.Wrap(err, "DeleteUser - Inserting user error")
	}

	return nil
}
