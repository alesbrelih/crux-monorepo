package grpc

import (
	"context"

	"github.com/alesbrelih/crux-monorepo/microservices/internal/grpc_clients"
	"github.com/alesbrelih/crux-monorepo/microservices/pkg"
	"github.com/alesbrelih/crux-monorepo/microservices/protos/build/services"
	"github.com/alesbrelih/crux-monorepo/microservices/services/registration/internal/repository"
	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	"github.com/trustelem/zxcvbn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewRegistrationService(log hclog.Logger, repo repository.Repository, userClient grpc_clients.UserClient) services.RegistrationServiceServer {
	return &registrationServiceServer{
		log:        log,
		repo:       repo,
		userClient: userClient,
	}
}

type registrationServiceServer struct {
	log        hclog.Logger
	repo       repository.Repository
	userClient grpc_clients.UserClient
	services.UnimplementedRegistrationServiceServer
}

func (s *registrationServiceServer) Register(ctx context.Context, request *services.RegisterRequest) (*services.RegisterResponse, error) {
	// whole this info could be put into JWT and passed as token beside invide
	// but this might be unnecesarry processing power?
	id := uuid.New().String()
	userInvite := repository.NewUserInviteFromRegistration(id, request)
	if err := s.repo.CreateInvite(userInvite); err != nil {
		s.log.Error("CreateInvite error. Error: %v", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	// TODO: send invite email

	type UserInvite struct {
		Id       string // this will be uid
		Name     string
		Surname  string
		Email    string
		Username string
	}

	return &services.RegisterResponse{
		Uuid: id,
	}, nil

}

func (s *registrationServiceServer) ConfirmRegistration(ctx context.Context, request *services.ConfirmRegistrationRequest) (*services.ConfirmRegistrationResponse, error) {
	userInvite, err := s.repo.GetInvite(request.Uuid)
	if err != nil {
		s.log.Error("ConfirmRegistration error. Error: %v", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	// TODO: write score lvl to some enviroment variable
	fields := pkg.GetStructStringValues(userInvite)
	if res := zxcvbn.PasswordStrength(request.Password, fields); res.Score < 3 {
		return nil, status.Error(codes.InvalidArgument, "Password strength is too low")
	}

	if request.Password != request.ConfirmPassword {
		return nil, status.Error(codes.InvalidArgument, "Passwords must match")
	}

	// call user microservice to create user
	id, err := s.userClient.CreateUser(ctx, userInvite.Username, userInvite.Email, request.Password)
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			// handled grpc error from user microservice -> forwarding it
			// maybe this can be simplyfied?
			s.log.Error("GRPC: ConfirmRegistration GRPC error calling CreateUser: %v", err.Error())
			return nil, status.Error(grpcErr.Code(), grpcErr.Message())
		}

		// unexpected server error
		s.log.Error("GRPC: ConfirmRegistration server error calling CreateUser: %v", err.Error())
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &services.ConfirmRegistrationResponse{
		Id: id,
	}, nil

}
