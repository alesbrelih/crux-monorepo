package grpc

import (
	"context"

	"github.com/alesbrelih/crux-monorepo/microservices/pkg"
	"github.com/alesbrelih/crux-monorepo/microservices/protos/build/services"
	"github.com/alesbrelih/crux-monorepo/microservices/services/user/internal/models"
	"github.com/alesbrelih/crux-monorepo/microservices/services/user/internal/repository"
	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewUserService(log hclog.Logger, passUtil pkg.PasswordUtil, repo repository.Repository) services.UserServiceServer {
	return &userServiceServer{
		log:      log,
		passUtil: passUtil,
		repo:     repo,
	}
}

type userServiceServer struct {
	log      hclog.Logger
	passUtil pkg.PasswordUtil
	repo     repository.Repository
	services.UnimplementedUserServiceServer
}

// accepts userinvite model with password set and writes to user db, Im using it this way because I want 1 microservices to hold
// 1 scope
func (s *userServiceServer) CreateUser(ctx context.Context, request *services.CreateUserRequest) (*services.CreateUserResponse, error) {

	// validation happens on confirm registration inside registration microservice...
	// all we need is to parse request and put it into database
	uuid, err := uuid.NewUUID()
	if err != nil {
		s.log.Error("GRPC: CreateUser - failed to create UUID", err.Error())
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	// move cost to some env
	hashed, err := s.passUtil.Hash([]byte(request.GetPassword()), 10)
	user := &models.User{
		Id:       uuid.String(),
		Email:    request.GetEmail(),
		Username: request.GetUsername(),
		Password: string(hashed),
	}

	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		s.log.Error("GRPC: CreateUser - failed to insert new user", err.Error())
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &services.CreateUserResponse{
		Id: id,
	}, nil

}

// this is used by other microservice that triggers create and if then fails
// it needs to call delete to revert action
func (s *userServiceServer) DeleteUser(ctx context.Context, request *services.DeleteUserRequest) (*services.DeleteUserResponse, error) {

	err := s.repo.DeleteUser(ctx, request.GetId())
	if err != nil {
		s.log.Error("GRPC: DeleteUser - failed to delete user with id: %s. Error: %s", request.GetId(), err.Error())
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &services.DeleteUserResponse{}, nil
}
