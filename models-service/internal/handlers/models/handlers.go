package models

import (
	"context"

	ssov3 "github.com/mgrankin-cloud/messenger/contract/gen/go/models"
	"github.com/mgrankin-cloud/messenger/internal/domain/models"
	"github.com/mgrankin-cloud/messenger/pkg/storage"
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
		userID int64,
	) (models.User, error)
	GetChat(
		ctx context.Context,
		chatID int64,
	) (models.Chat, error)
	GetApp(
		ctx context.Context,
		appID int64,
	) (models.App, error)
	GetMessage(
		ctx context.Context,
		messageID int64,
	) (models.Message, error)
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
		Photo:    user.Photo,
		Active:   user.Active,
	}, nil
}

func (s *serverAPI) GetChat(
	ctx context.Context,
	request *ssov3.GetChatRequest,
) (*ssov3.GetChatResponse, error) {

	if request.GetChatId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "message id is required")
	}

	chat, err := s.storage.GetChat(ctx, request.GetChatId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "message not found")
	}

	return &ssov3.GetChatResponse{
		ChatId:       chat.ID,
		ChatName:     chat.Name,
		Photo:        chat.Photo,
		ChatType:     ssov3.ChatType(chat.ChatType),
		Participants: make([]*ssov3.GetUserResponse, 0),
	}, nil
}

func (s *serverAPI) GetMessage(
	ctx context.Context,
	request *ssov3.GetMessageRequest,
) (*ssov3.GetMessageResponse, error) {

	if request.GetMessageId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "message id is required")
	}

	message, err := s.storage.Message(ctx, request.GetMessageId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "message not found")
	}

	return &ssov3.GetMessageResponse{
		MessageId:  message.ID,
		Content:    message.Content,
		CreatedAt:  message.CreatedAt,
		CreatedBy:  message.CreatedBy,
		ReplyTo:    message.ReplyToID,
		ReceivedBy: message.ReceivedBy,
		ReceivedAt: message.ReceivedAt,
		IsRead:     message.IsRead,
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
		AppId:  int64(app.ID),
		Name:   app.Name,
		Secret: app.Secret,
	}, nil
}
