package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/alesbrelih/crux-monorepo/microservices/internal/repository_connect"
	"github.com/alesbrelih/crux-monorepo/microservices/services/mail/internal/models"
	"github.com/pkg/errors"
)

func NewRepository(dsn string) Repository {
	return &repositoryPQ{
		conn: repository_connect.NewRepositoryConnect(dsn),
	}
}

type Repository interface {
	Get(ctx context.Context, id int64) (*models.Mail, error)
	GetAll(ctx context.Context, from, to time.Time, reciever, status string) ([]*models.Mail, error)
	ToQueue(ctx context.Context, reciever, subject, body string) (int64, error)
}

type repositoryPQ struct {
	conn repository_connect.RepositoryConnect
}

// gets single mail details from db
func (r *repositoryPQ) Get(ctx context.Context, id int64) (*models.Mail, error) {
	conn, err := r.conn.Connect()
	if err != nil {
		return nil, errors.Wrap(err, "Get: error connecting to db.")
	}
	defer conn.Close()

	q := "SELECT id, to, subject, body, status, created_at WHERE id := $1"

	var mail models.Mail
	err = conn.QueryRowContext(ctx, q, id).Scan(&mail.Id, &mail.Reciever, &mail.Subject, &mail.Body, &mail.Status, &mail.CreatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "Get: error retrieving Mail")
	}
	return &mail, nil
}

// gets all mail from database
func (r *repositoryPQ) GetAll(ctx context.Context, from, to time.Time, reciever, status string) ([]*models.Mail, error) {
	conn, err := r.conn.Connect()
	if err != nil {
		return nil, errors.Wrap(err, "GetAll: error connecting to db.")
	}
	defer conn.Close()

	// build query
	q := "SELECT id, reciever, status, created_at FROM mail WHERE 1"
	// https://play.golang.org/p/BlADhht9PwO
	params := make([]interface{}, 4)
	if !from.IsZero() {
		q += " AND created_at >= $" + strconv.Itoa(len(params)+1)
		params = append(params, from)
	}
	if !to.IsZero() {
		q += " AND created_at <= $" + strconv.Itoa(len(params)+1)
		params = append(params, to)
	}
	if reciever != "" {
		q += " AND reciever LIKE '%$" + strconv.Itoa(len(params)+1) + "%'"
		params = append(params, reciever)
	}
	if status != "" {
		q += " AND status = $" + strconv.Itoa(len(params)+1)
		params = append(params, status)
	}

	rows, err := conn.QueryContext(ctx, q, params...)
	if err != nil {
		return nil, errors.Wrap(err, "GetAll: error retrieving mails from db.")
	}

	mails := make([]*models.Mail, 0)
	for rows.Next() {
		mail := new(models.Mail)
		if err := rows.Scan(&mail.Id, &mail.Reciever, &mail.Status, &mail.CreatedAt); err != nil {
			return nil, errors.Wrap(err, "GetAll: error scanning mail")
		}
		mails = append(mails, mail)
	}
	return mails, nil

}

// Inserts mail into queue
func (r *repositoryPQ) ToQueue(ctx context.Context, reciever, subject, body string) (int64, error) {
	conn, err := r.conn.Connect()
	if err != nil {
		return 0, errors.Wrap(err, "ToQueue: error connecting to db.")
	}
	defer conn.Close()

	q := `INSERT INTO crux_mail (reciever, subject, body, status, created_at)
				VALUES ($1, $2, $3, $4, (now() at time zone 'utc')) 
			RETURNING id`
	var id int64
	err = conn.QueryRowContext(ctx, q, reciever, subject, body, models.DEFAULT_MAIL_STATUS).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "ToQueue: error inserting to db")
	}
	return id, nil
}
