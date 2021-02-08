package repository

import (
	"context"
	"database/sql"
	"fmt"

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
	BefriendUser(ctx context.Context, userOne, userTwo string) error
	UnfriendUser(ctx context.Context, userOne, userTwo string) error
	HandleFriendInvite(ctx context.Context, userOne, userTwo string, isAccepted bool) error
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

// fromUser sends invitation to toUser
func (repo *repositoryPQ) BefriendUser(ctx context.Context, userOne, userTwo string) error {
	conn, err := repo.conn.Connect()

	if err != nil {
		return errors.Wrap(err, "BefriendUser - DB connecting error")
	}

	defer conn.Close()

	q := "INSERT INTO social_network (user_one, user_two) VALUES ($1, $2)"
	_, err = conn.ExecContext(ctx, q, userOne, userTwo)
	if err != nil {
		return errors.Wrap(err, "BefriendUser - DB insert error")
	}

	return nil

}

func (repo *repositoryPQ) UnfriendUser(ctx context.Context, userOne, userTwo string) error {
	conn, err := repo.conn.Connect()
	if err != nil {
		return errors.Wrap(err, "UnfriendUser - DB connecting error")
	}
	defer conn.Close()

	q := "DELETE FROM social_network WHERE user_one IN ($1, $2) AND user_two IN ($1, $2)"
	_, err = conn.ExecContext(ctx, q, userOne, userTwo)
	if err != nil {
		return errors.Wrap(err, "UnfriendUser - DB remove error")
	}

	return nil
}

func (repo *repositoryPQ) HandleFriendInvite(ctx context.Context, userOne, userTwo string, isAccepted bool) error {
	conn, err := repo.conn.Connect()
	if err != nil {
		return errors.Wrap(err, "HandleFriendInvite - DB connecting error")
	}
	defer conn.Close()

	var id int64
	q := "SELECT id FROM social_network WHERE user_one IN ($1, $2) AND user_two IN ($1, $2)"
	err = conn.QueryRowContext(ctx, q, userTwo).Scan(&id)
	if err != nil {
		return errors.Wrap(err, "HandleFriendInvite - DB query social network error")
	}

	// TODO: rework this -> using different table for invites!
	// maybe? not sure
	if isAccepted {
		return repo.acceptFriendInvite(ctx, conn, id)
	}
	return repo.declineFriendInvite(ctx, conn, id)

}

func (repo *repositoryPQ) acceptFriendInvite(ctx context.Context, conn *sql.DB, id int64) error {
	q := "UPDATE social_network SET social_status = 'ACCEPTED', date_accepted = (NOW() at time zone 'utc') WHERE id  = $1"
	res, err := conn.ExecContext(ctx, q, id)
	if err != nil {
		return errors.Wrap(err, "acceptFriendInvite - error saving accept to db")
	}

	if num, err := res.RowsAffected(); err != nil {
		if err != nil {
			return errors.Wrap(err, "acceptFriendInvite - error retrieving rows affected")
		}
		if num == 0 {
			return errors.Wrap(fmt.Errorf("No such invitation id: %v", id), "acceptFriendInvite - error saving accept to db")
		}
	}

	return nil

}

func (repo *repositoryPQ) declineFriendInvite(ctx context.Context, conn *sql.DB, id int64) error {
	q := "DELETE FROM social_network WHERE id  = $1"
	res, err := conn.ExecContext(ctx, q, id)
	if err != nil {
		return errors.Wrap(err, "declineFriendInvite - error saving decline to db")
	}

	if num, err := res.RowsAffected(); err != nil {
		if err != nil {
			return errors.Wrap(err, "declineFriendInvite - error retrieving rows affected")
		}
		if num == 0 {
			return errors.Wrap(fmt.Errorf("No such invitation id: %v", id), "declineFriendInvite - error saving accept to db")
		}
	}

	return nil
}
