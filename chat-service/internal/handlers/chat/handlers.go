package user

import (
	"context"
	"errors"

	"log/slog"

	ssov4 "github.com/mgrankin-cloud/messenger/contract/gen/go/chat"
	"github.com/mgrankin-cloud/messenger/internal/domain/models"
	errors2 "github.com/mgrankin-cloud/messenger/pkg/storage"

	"github.com/mgrankin-cloud/messenger/pkg/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	ssov4.UnimplementedChatServer
	chat    Chat
	storage *storage.Storage
	log     *slog.Logger
}

type Chat interface {
	GetChat(
		ctx context.Context,
		chatID int64,
	) (models.Chat, error)
	UpdateChat(
		ctx context.Context,
		chatName string,
		photo []byte,
		chatID int64,
	) (success bool, err error)
	CreateChat(
		ctx context.Context,
		name string,
		photo []byte,
		chatType ssov4.ChatType,
	) (chatID int64, err error)
	DeleteChat(
		ctx context.Context,
		chatID int64,
	) (err error)
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

	chat, err := s.chat.GetChat(ctx, request.GetChatId())
	if err != nil {
		if errors.Is(err, errors2.ErrChatNotFound) {
			return nil, status.Error(codes.NotFound, "chat not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get chat: %v", err)
	}

	return &ssov4.GetChatResponse{	
		ChatName: chat.Name,
		ChatType: ssov4.ChatType(chat.ChatType),
		Photo:    chat.Photo,
	}, nil
}

func (s *serverAPI) CreateChat(
	ctx context.Context,
	request *ssov4.CreateChatRequest,
) (*ssov4.CreateChatResponse, error) {
	if request.GetChatName() == "" {
		return nil, status.Error(codes.InvalidArgument, "chat name is required")
	}

	if request.GetChatType() == 0 {
		return nil, status.Error(codes.InvalidArgument, "chat type is required")
	}

	chatID, err := s.chat.CreateChat(ctx, request.GetChatName(), request.GetPhoto(), request.GetChatType())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create chat: %v", err)
	}

	return &ssov4.CreateChatResponse{
		Success: true,
		ChatId:  chatID,
	}, nil
}

func (s *serverAPI) UpdateChat(
	ctx context.Context,
	request *ssov4.UpdateChatRequest,
) (*ssov4.UpdateChatResponse, error) {
	if request.GetChatId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "chat id is required")
	}

	if request.GetChatName() == "" {
		return nil, status.Error(codes.InvalidArgument, "chat name is required")
	}

	success, err := s.chat.UpdateChat(ctx, request.GetChatName(), request.GetPhoto(), request.GetChatId())
	if err != nil {
		if errors.Is(err, errors2.ErrChatNotFound) {
			return nil, status.Errorf(codes.NotFound, "chat not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update chat: %v", err)
	}

	return &ssov4.UpdateChatResponse{
		Success: success,
	}, nil
}

func (s *serverAPI) DeleteChat(
	ctx context.Context,
	request *ssov4.DeleteChatRequest,
) (*ssov4.DeleteChatResponse, error) {
	if request.GetChatId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "chat id is required")
	}

	err := s.chat.DeleteChat(ctx, request.GetChatId())
	if err != nil {
		if errors.Is(err, storage.ErrChatNotFound) {
			return nil, status.Errorf(codes.NotFound, "chat not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete chat: %v", err)
	}

	return &ssov4.DeleteChatResponse{
		Success: true,
	}, nil
}

func (s *serverAPI) logAction(userID, chatID int64, action string) {
	s.log.Info("User action", slog.Int64("userID", userID), slog.Int64("chatID", chatID), slog.String("action", action))
}
