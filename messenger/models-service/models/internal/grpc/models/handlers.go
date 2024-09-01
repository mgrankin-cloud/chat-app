package models

import (
	"Messenger-android/messenger/internal/storage"
	ssov3 "Messenger-android/messenger/models-service/contract/gen/go/models"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	ssov3.UnimplementedModelsServer
	models  Models
	storage *storage.Storage
}

type Models interface {
	GetUser(
		ctx context.Context,
		request *ssov3.GetUserRequest,
	) (*ssov3.GetUserResponse, error)
	GetChat(
		ctx context.Context,
		request *ssov3.GetChatRequest,
	) (*ssov3.GetChatResponse, error)
	GetApp(
		ctx context.Context,
		request *ssov3.GetAppRequest,
	) (*ssov3.GetAppResponse, error)
	GetMessage(
		ctx context.Context,
		request *ssov3.GetMessageRequest,
	) (*ssov3.GetMessageResponse, error)
}

func RegisterModelsService(gRPCServer *grpc.Server, models Models) {
	ssov3.RegisterModelsServer(gRPCServer, &serverAPI{models: models})
}

func (s *serverAPI) GetUser(
	ctx context.Context,
	request *ssov3.GetUserRequest,
) (*ssov3.GetUserResponse, error) {

	if request.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	user, err := s.storage.GetUserByID(ctx, request.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &ssov3.GetUserResponse{
		UserId:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		Phone:    user.Phone,
	}, nil
}

func (s *serverAPI) GetChat(
	ctx context.Context,
	request *ssov3.GetChatRequest,
) (*ssov3.GetChatResponse, error) {

	if request.GetChatId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "message id is required")
	}

	chat, err := s.storage.GetChatByID(ctx, request.GetChatId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "message not found")
	}

	return &ssov3.GetChatResponse{
		ChatId:       chat.ID,
		ChatName:     chat.Name,
		Photo:        chat.Photo,
		Participants: chat.Participants,
		ChatType:     chat.ChatType,
	}, nil
}

func (s *serverAPI) GetMessage(
	ctx context.Context,
	request *ssov3.GetMessageRequest,
) (*ssov3.GetMessageResponse, error) {

	if request.GetMessageId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "message id is required")
	}

	chat, err := s.storage.GetMessageByID(ctx, request.GetMessageId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "message not found")
	}

	return &ssov3.GetMessageResponse{
		MessageId:  message.ID,
		ChatID:     chat.ID,
		Content:    message.Name,
		CreatedAt:  message.CreatedAt,
		CreatedBy:  message.CreatedBy,
		ReplyTo:    message.ReplyTo,
		ReceivedBy: message.ReceivedBy,
		ReceivedAt: message.ReceivedAt,
	}, nil
}

func (s *serverAPI) GetApp(
	ctx context.Context,
	request *ssov3.GetAppRequest,
) (*ssov3.GetAppResponse, error) {

	if request.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "app id is required")
	}

	app, err := s.storage.GetAppByID(ctx, request.GetAppId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "app not found")
	}

	return &ssov3.GetAppResponse{
		AppId:  app.AppID,
		Name:   app.Name,
		Secret: app.Secret,
	}, nil
}
