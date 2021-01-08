package grpc

import (
	"context"

	"github.com/alesbrelih/crux-monorepo/microservices/pkg"
	"github.com/alesbrelih/crux-monorepo/microservices/protos/build/services"
	"github.com/alesbrelih/crux-monorepo/microservices/services/registration/internal/repository"
	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	"github.com/trustelem/zxcvbn"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewRegistrationService(log hclog.Logger, repo repository.Repository) services.RegistrationServiceServer {
	return &registrationServiceServer{
		log:  log,
		repo: repo,
	}
}

type registrationServiceServer struct {
	log  hclog.Logger
	repo repository.Repository
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

	// TODO: refactor in something unittestable
	conn, err := grpc.Dial("")
	if err != nil {
		s.log.Error("GRPC: ConfirmRegistration error connecting to user gprc. Error: %v", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	defer conn.Close()

	createUserRequest := &services.CreateUserRequest{
		Username: userInvite.Username,
		Email:    userInvite.Email,
		Password: request.Password,
	}
	c := services.NewUserServiceClient(conn)
	resp, err := c.CreateUser(ctx, createUserRequest)
	if err != nil {
		s.log.Error("GRPC: ConfirmRegistration, error from user.CreateUser. Error: %s", err.Error())
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &services.ConfirmRegistrationResponse{
		Id: resp.GetId(),
	}, nil
}
