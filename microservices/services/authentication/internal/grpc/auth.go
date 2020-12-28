package grpc

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/alesbrelih/crux-monorepo/microservices/protos/build/services"
	"github.com/alesbrelih/crux-monorepo/microservices/services/authentication/internal/my_jwt"
	"github.com/alesbrelih/crux-monorepo/microservices/services/authentication/internal/repository"
	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewAuthService(jwt my_jwt.JwtService, repository repository.Repository) services.AuthServiceServer {
	return &authService{
		repository: repository,
		jwt:        jwt,
	}
}

type authService struct {
	log        hclog.Logger
	jwt        my_jwt.JwtService
	repository repository.Repository
	services.UnimplementedAuthServiceServer
}

func (auth *authService) Authenticate(ctx context.Context, request *services.AuthenticateRequest) (*services.AuthenticateResponse, error) {

	// get from db if exists
	user, err := auth.repository.GetUserByUsername(ctx, request.GetUsername())
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, status.Error(codes.Unauthenticated, "Invalid username or password")
		}
		auth.log.Error("Error occured in authService.Authenticate. Error: %#v", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	// validate password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.GetPassword()))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid username or password")
	}

	tokenPair, err := auth.jwt.GenerateJwtPair(strconv.FormatInt(user.Id, 10))
	if err != nil {
		auth.log.Error("Error occured in authService.Authenticate when Generating JWT pair. Error: %#v", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &services.AuthenticateResponse{
		Access:  tokenPair.Access,
		Refresh: tokenPair.Refresh,
	}, nil
}

func (auth *authService) RefreshToken(ctx context.Context, request *services.RefreshTokenRequest) (*services.RefreshTokenResponse, error) {
	claims, err := auth.jwt.GetClaims(request.GetRefresh())
	if err != nil {
		if errors.Cause(err) == my_jwt.InvalidTokenError {
			auth.log.Error("Error occured in RefreshToken. Error: %v", err)
			return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
		}
		auth.log.Error("Error occured in RefreshToken. Error: %v", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	id, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		auth.log.Error("Error occured in RefreshToken. Error: %v", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	if active, err := auth.repository.HasAccess(ctx, int64(id)); err != nil || !active {
		if err != nil {
			auth.log.Error("Error occured in RefreshToken. Error: %v", err)
			return nil, status.Error(codes.Internal, "Internal server error")
		}
		if !active {
			return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
		}
	}

	tokenPair, err := auth.jwt.GenerateJwtPair(claims.Subject)
	if err != nil {
		auth.log.Error("Error generating jwt pair for id: %v. Error: %v", id, err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &services.RefreshTokenResponse{
		Access:  tokenPair.Access,
		Refresh: tokenPair.Refresh,
	}, nil

}
