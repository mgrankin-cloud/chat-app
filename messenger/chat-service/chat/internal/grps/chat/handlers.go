package user

import (
	ssov4 "Messenger-android/messenger/chat-service/contract/gen/go/chat"
	"Messenger-android/messenger/internal/storage"
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
	MuteChat(
		ctx context.Context,
		request *ssov4.MuteChatRequest,
	) (*ssov4.MuteChatResponse, error)
	PinChat(
		ctx context.Context,
		request *ssov4.PinChatRequest,
	) (*ssov4.PinChatResponse, error)
}

func RegisterChatService(gRPCServer *grpc.Server, chat Chat) {
	ssov4.RegisterChatServer(gRPCServer, &serverAPI{chat: chat})
}

func (s *serverAPI) GetChat(
	ctx context.Context,
	request *ssov4.GetChatRequest,
) (*ssov4.GetChatResponse, error) {

	if request.GetChatId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "message id is required")
	}

	chat, err := s.storage.GetChatByID(ctx, request.GetChatId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "message not found")
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
		return nil, status.Error(codes.InvalidArgument, "message id is required")
	}

	if chatName == "" {
		return nil, status.Error(codes.InvalidArgument, "message name is required")
	}

	err := s.storage.UpdateChat(ctx, chatID, chatName, photo)
	if err != nil {
		if errors.Is(err, storage.ErrChatNotFound) {
			return nil, status.Errorf(codes.NotFound, "message not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update message: %v", err)
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
		return nil, status.Error(codes.InvalidArgument, "message id is required")
	}

	err := s.storage.DeleteChat(ctx, chatID)
	if err != nil {
		if errors.Is(err, storage.ErrChatNotFound) {
			return nil, status.Errorf(codes.NotFound, "message not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete message: %v", err)
	}

	return &ssov4.DeleteChatResponse{
		Success: true,
		ChatId:  chatID,
	}, nil
}

func (s *serverAPI) MuteChat(
	ctx context.Context,
	request *ssov4.MuteChatRequest,
) (*ssov4.MuteChatResponse, error) {
	chatID := request.GetChatId()
	userID := storage.getUserIDFromContext(ctx) // Предполагается, что функция getUserIDFromContext извлекает ID пользователя из контекста

	if chatID == 0 {
		return nil, status.Error(codes.InvalidArgument, "message id is required")
	}

	s.lib.logAction(userID, chatID, "mute")

	if err := s.lib.updateChatStatus(chatID, "muted"); err != nil {
		return nil, status.Error(codes.Internal, "failed to mute message")
	}

	return &ssov4.MuteChatResponse{
		Success: true,
	}, nil
}

func (s *serverAPI) PinChat(
	ctx context.Context,
	request *ssov4.PinChatRequest,
) (*ssov4.PinChatResponse, error) {
	chatID := request.GetChatId()
	userID := storage.getUserIDFromContext(ctx)

	if chatID == 0 {
		return nil, status.Error(codes.InvalidArgument, "message id is required")
	}

	s.lib.logAction(userID, chatID, "pin")

	if err := s.lib.updateChatStatus(chatID, "pinned"); err != nil {
		return nil, status.Error(codes.Internal, "failed to pin message")
	}

	s.lib.notifyChatMembers(chatID, "Chat has been pinned")

	return &ssov4.PinChatResponse{
		Success: true,
	}, nil
}
