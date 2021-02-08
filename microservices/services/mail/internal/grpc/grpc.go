package grpc

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/alesbrelih/crux-monorepo/microservices/pkg"
	"github.com/alesbrelih/crux-monorepo/microservices/protos/build/services"
	"github.com/alesbrelih/crux-monorepo/microservices/services/mail/internal/models"
	"github.com/alesbrelih/crux-monorepo/microservices/services/mail/internal/repository"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewMailService(log hclog.Logger, repo repository.Repository) services.MailServiceServer {
	return &mailServiceServer{
		log:  log,
		repo: repo,
	}
}

type mailServiceServer struct {
	log  hclog.Logger
	repo repository.Repository
	services.UnimplementedMailServiceServer
}

// Sends mail -> puts it in db queue
func (s *mailServiceServer) SendMail(ctx context.Context, request *services.SendMailRequest) (*services.SendMailResponse, error) {

	// even though we validate email thought grpc validators
	// that validation is basic, here is better one with MX lookup aswell
	if !pkg.IsEmailValid(request.Reciever) {
		return nil, status.Error(codes.InvalidArgument, "Invalid email address")
	}

	id, err := s.repo.ToQueue(ctx, request.GetReciever(), request.GetSubject(), request.GetBody())
	if err != nil {
		s.log.Error(
			"SendMail: error saving mail to queue. Reciever: '%s', Subject: '%s', Body: '%s'. Error: '%s'",
			request.GetReciever(),
			request.GetSubject(),
			request.GetBody(),
			err.Error(),
		)
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	return &services.SendMailResponse{
		Id: id,
	}, nil
}

// gets single mail
func (s *mailServiceServer) GetMail(ctx context.Context, request *services.GetMailRequest) (*services.GetMailResponse, error) {
	mail, err := s.repo.Get(ctx, request.GetId())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "Mail not found")
		}
		s.log.Error("GetMail: error retrieving mail from db. Id: '%d'. Error: '%s'", request.GetId(), err.Error())
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return s.toGetMailResponse(mail)

}

// gets all mails
func (s *mailServiceServer) GetMails(ctx context.Context, request *services.GetMailsRequest) (*services.GetMailsResponse, error) {
	mails, err := s.repo.GetAll(ctx, request.From.AsTime(), request.To.AsTime(), request.GetReciever(), request.GetStatus().String())
	if err != nil {
		s.log.Error(
			"GetMails: error querying database.  From: '%v', To: '%v', Reciever: '%s', Status: '%s'. Error: '%s'",
			request.From.AsTime(),
			request.To.AsTime(),
			request.GetReciever(),
			request.GetStatus().String(),
		)
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	return s.toGetMailsResponse(mails)
}

// transforms to grpc GetMailResponse for single mail item
func (s *mailServiceServer) toGetMailResponse(mail *models.Mail) (*services.GetMailResponse, error) {

	mailStatus, err := s.getEnumFromDb(mail.Status)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	return &services.GetMailResponse{
		Id:        mail.Id,
		Reciever:  mail.Reciever,
		CreatedAt: timestamppb.New(mail.CreatedAt),
		Status:    mailStatus,
	}, nil
}

// transforms to grpc GetMailsResponse
func (s *mailServiceServer) toGetMailsResponse(mails []*models.Mail) (*services.GetMailsResponse, error) {

	items := make([]*services.GetMailsItem, len(mails))

	for _, mail := range mails {
		mailStatus, err := s.getEnumFromDb(mail.Status)
		if err != nil {
			return nil, status.Error(codes.Internal, "Internal server error")
		}

		item := &services.GetMailsItem{
			Id:        mail.Id,
			Reciever:  mail.Reciever,
			CreatedAt: timestamppb.New(mail.CreatedAt),
			Status:    mailStatus,
		}
		items = append(items, item)
	}
	return &services.GetMailsResponse{
		Items: items,
	}, nil
}

func (s *mailServiceServer) getEnumFromDb(dbValue string) (services.MailStatus, error) {
	grpcEnumVal, ok := services.MailStatus_value[dbValue]
	if !ok {
		// setting to -1 since its int value and it has 0 values in grpc implementation
		s.log.Error("GetMail: invalid enum type for status. Got: %s", dbValue)
		return -1, fmt.Errorf("Invalid enum type for status. Got: %s", dbValue)
	}
	return services.MailStatus(grpcEnumVal), nil
	// interesting:
	// type MailStatus int32

	// const (
	// 	MailStatus_IN_QUEUE MailStatus = 0
	// 	MailStatus_SENT     MailStatus = 1
	// 	MailStatus_ERROR    MailStatus = 2
	// )
	// services.MailStatus(val)
}
