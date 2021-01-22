package grpc_clients

import (
	"context"

	"github.com/alesbrelih/crux-monorepo/microservices/pkg/grpc_helper"
	"github.com/alesbrelih/crux-monorepo/microservices/protos/build/services"
	"google.golang.org/grpc"
)

func NewUserClient(address string) UserClient {
	return &userClient{
		grpc_helper.GRPCRequest{
			Address: address,
		},
	}
}

type UserClient interface {
	// calls user service and creates user with specified username, email, password
	// and returns id
	CreateUser(ctx context.Context, username, email, password string) (string, error)
}

type userClient struct {
	grpc_helper.GRPCRequest
}

func (u *userClient) CreateUser(ctx context.Context, username, email, password string) (string, error) {
	result, err := u.CallInsecureBlocking(func(conn *grpc.ClientConn) (interface{}, error) {
		client := services.NewUserServiceClient(conn)

		req := &services.CreateUserRequest{
			Username: username,
			Email:    email,
			Password: password,
		}

		res, err := client.CreateUser(ctx, req)
		if err != nil {
			return "", nil

		}
		return res.GetId(), nil
	})

	// no need to check assertion
	// since func is above this
	return result.(string), err
}
