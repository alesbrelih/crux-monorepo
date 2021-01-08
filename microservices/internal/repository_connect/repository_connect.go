package repository_connect

import "database/sql"

func NewRepositoryConnect(dsn string) RepositoryConnect {
	return &repositoryConnectPQ{
		dsn: dsn,
	}
}

type RepositoryConnect interface {
	Connect() (*sql.DB, error)
	GetDsn() string
}

type repositoryConnectPQ struct {
	dsn string
}

func (d *repositoryConnectPQ) Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", d.dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
func (d *repositoryConnectPQ) GetDsn() string {
	return d.dsn
}
