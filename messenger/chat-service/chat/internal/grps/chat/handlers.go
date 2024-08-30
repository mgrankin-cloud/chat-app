package user

import (
	ssov4 "Messenger-android/messenger/chat-service/contract/gen/go/chat"
	"Messenger-android/messenger/logger/validator"
	"Messenger-android/messenger/storage"
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	ssov4.UnimplementedChatServer
	chat    Chat
	storage *storage.Storage
}

type Chat interface {
	GetChat(
		ctx context.Context,
		request *ssov4.GetChatRequest,
	) (*ssov4.GetChatResponse, error)
	UpdateChat(
		ctx context.Context,
		request *ssov4.UpdateChatRequest,
	) (*ssov4.UpdateChatResponse, error)
	DeleteChat(
		ctx context.Context,
		request *ssov4.DeleteChatRequest,
	) (*ssov4.DeleteChatResponse, error)
}

func RegisterChatService(gRPCServer *grpc.Server, chat Chat) {
	ssov4.RegisterChatServer(gRPCServer, &serverAPI{chat: chat})
}

func (s *serverAPI) GetChat(
	ctx context.Context,
	request *ssov4.GetChatRequest,
) (*ssov4.GetChatResponse, error) {

	if request.GetChatId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "chat id is required")
	}

	chat, err := s.storage.GetChatByID(ctx, request.GetChatId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "chat not found")
	}

	return &ssov4.GetChatResponse{
		ChatId:   chat.ID,
		ChatName: chat.Name,
	}, nil
}

func (s *serverAPI) UpdateChat(
	ctx context.Context,
	request *ssov4.UpdateChatRequest,
) (*ssov4.UpdateChatResponse, error) {
	chatID := request.GetChatId()
	chatName := request.GetChatName()
	photo := request.GetPhoto()

	if chatID == 0 {
		return nil, status.Error(codes.InvalidArgument, "chat id is required")
	}

	if err := validator.ValidateField(request.Email, "chat_name", validator.IsValidUsername); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.storage.UpdateChat(ctx, chatID, chatName, photo)
	if err != nil {
		if errors.Is(err, storage.ErrChatNotFound) {
			return nil, status.Errorf(codes.NotFound, "chat not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update chat: %v", err)
	}

	return &ssov4.UpdateChatResponse{
		Success: true,
		ChatId:  chatID,
	}, nil
}

func (s *serverAPI) DeleteChat(
	ctx context.Context,
	request *ssov4.DeleteChatRequest,
) (*ssov4.DeleteChatResponse, error) {
	chatID := request.GetChatId()

	if chatID == 0 {
		return nil, status.Error(codes.InvalidArgument, "chat id is required")
	}

	err := s.storage.DeleteChat(ctx, cahtID)
	if err != nil {
		if errors.Is(err, storage.ErrChatNotFound) {
			return nil, status.Errorf(codes.NotFound, "chat not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete chat: %v", err)
	}

	return &ssov4.DeleteChatResponse{
		Success: true,
		ChatId:  chatID,
	}, nil
}
