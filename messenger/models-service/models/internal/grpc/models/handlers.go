package models

import (
	ssov3 "Messenger-android/messenger/models-service/contract/gen/go/models"
	"Messenger-android/messenger/storage"
	ssov2 "Messenger-android/messenger/user-service/contract/gen/go/user"
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
	) (*ssov2.GetUserResponse, error)
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
		return nil, status.Error(codes.InvalidArgument, "chat id is required")
	}

	chat, err := s.storage.GetChatByID(ctx, request.GetChatId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "chat not found")
	}

	return &ssov3.GetChatResponse{
		ChatId:       chat.ID,
		ChatName:     chat.Name,
		PhotoUrl:     chat.PhotoUrl,
		Participants: chat.Participants,
		ChatType:     chat.ChatType,
	}, nil
}

func (s *serverAPI) GetMessage(
	ctx context.Context,
	request *ssov3.GetMessageRequest,
) (*ssov3.GetMessageResponse, error) {

	if request.GetMessageId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "chat id is required")
	}

	chat, err := s.storage.GetMessageByID(ctx, request.GetMessageId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "chat not found")
	}

	return &ssov3.GetMessageResponse{
		ChatId:       chat.ID,
		ChatName:     chat.Name,
		PhotoUrl:     chat.PhotoUrl,
		Participants: chat.Participants,
		ChatType:     chat.ChatType,
	}, nil
}

func (s *serverAPI) GetApp(
	ctx context.Context,
	request *ssov3.GetChatRequest,
) (*ssov3.GetChatResponse, error) {

	if request.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "chat id is required")
	}

	chat, err := s.storage.GetAppByID(ctx, request.GetAppId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "chat not found")
	}

	return &ssov3.GetAppResponse{
		ChatId:       chat.ID,
		ChatName:     chat.Name,
		PhotoUrl:     chat.PhotoUrl,
		Participants: chat.Participants,
		ChatType:     chat.ChatType,
	}, nil
}
