package repository

import (
	"github.com/alesbrelih/crux-monorepo/microservices/internal/repository_connect"
	"github.com/alesbrelih/crux-monorepo/microservices/protos/build/services"
)

type UserInvite struct {
	Id       string // this will be uid
	Name     string
	Surname  string
	Email    string
	Username string
}

func NewUserInviteFromRegistration(id string, registration *services.RegisterRequest) *UserInvite {
	return &UserInvite{
		Id:       id,
		Email:    registration.GetEmail(),
		Username: registration.GetUsername(),
	}
}

func NewRepository(dsn string) Repository {
	return &repositoryPQ{
		conn: repository_connect.NewRepositoryConnect(dsn),
	}
}

type Repository interface {
	CreateInvite(userInvite *UserInvite) error
	GetInvite(id string) (*UserInvite, error)
	IsValid(token string) error
	DeleteInvite(email string) error
	DeleteOverdue() error
}

type repositoryPQ struct {
	conn repository_connect.RepositoryConnect
}

func (repo *repositoryPQ) CreateInvite(userInvite *UserInvite) error {
	panic("NOT YET IMPLEMENTED")
}

func (repo *repositoryPQ) GetInvite(id string) (*UserInvite, error) {
	panic("NOT YET IMPLEMENTED")
}

func (repo *repositoryPQ) IsValid(token string) error {
	panic("NOT YET IMPLEMENTED")
}

func (repo *repositoryPQ) DeleteInvite(email string) error {
	panic("NOT YET IMPLEMENTED")
}

func (repo *repositoryPQ) DeleteOverdue() error {
	panic("NOT YET IMPLEMENTED")
}
